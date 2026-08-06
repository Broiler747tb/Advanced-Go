package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"order-api/Db"
	"order-api/Order"
	"order-api/Product"
	"order-api/User"
	"order-api/pkg/jwt"
)

const (
	testSecret = "e2e-secret"
	testPhone  = 89990009901
)

type testEnv struct {
	Db       *Db.Db
	Router   *http.ServeMux
	User     *User.User
	Products []Product.Product
	Token    string
}

func initData(t *testing.T) *testEnv {
	t.Helper()

	_ = godotenv.Load("../.env")
	dsn := os.Getenv("DSN_TEST")
	if dsn == "" {
		t.Skip("DSN_TEST is not set, skipping e2e test")
	}

	connection, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Skipf("test database is not reachable: %v", err)
	}
	database := &Db.Db{DB: connection}

	if err := database.AutoMigrate(&Product.Product{}, &User.User{}, &Order.Order{}); err != nil {
		t.Fatal(err)
	}

	database.DB.Unscoped().Where("phone = ?", testPhone).Delete(&User.User{})

	user := &User.User{Phone: testPhone, Name: "e2e"}
	if err := database.DB.Create(user).Error; err != nil {
		t.Fatal(err)
	}

	var products []Product.Product
	for _, name := range []string{"e2e product a", "e2e product b", "e2e product c"} {
		product := Product.NewProduct(name)
		if err := database.DB.Create(product).Error; err != nil {
			t.Fatal(err)
		}
		products = append(products, *product)
	}

	token, err := jwt.NewJWT(testSecret).Create(jwt.JWTData{Phone: user.Phone})
	if err != nil {
		t.Fatal(err)
	}

	router := http.NewServeMux()
	NewOrderHandler(router, OrderHandlerDeps{
		OrderRepository:   Order.NewOrderRepository(database),
		ProductRepository: Product.NewProductRepository(database),
		UserRepository:    User.NewUserRepository(database),
		Secret:            testSecret,
	})

	return &testEnv{
		Db:       database,
		Router:   router,
		User:     user,
		Products: products,
		Token:    token,
	}
}

func (env *testEnv) removeData(t *testing.T) {
	t.Helper()

	database := env.Db.DB

	var orders []Order.Order
	database.Unscoped().Where("user_id = ?", env.User.ID).Find(&orders)
	for _, order := range orders {
		if err := database.Model(&order).Association("Products").Clear(); err != nil {
			t.Error(err)
		}
	}
	database.Unscoped().Where("user_id = ?", env.User.ID).Delete(&Order.Order{})

	for _, product := range env.Products {
		database.Unscoped().Delete(&Product.Product{}, product.ID)
	}
	database.Unscoped().Delete(&User.User{}, env.User.ID)
}

func TestCreateOrder(t *testing.T) {
	env := initData(t)
	defer env.removeData(t)

	productIds := []uint{env.Products[0].ID, env.Products[1].ID}
	body, err := json.Marshal(Order.OrderCreateRequest{ProductIds: productIds})
	if err != nil {
		t.Fatal(err)
	}

	request := httptest.NewRequest(http.MethodPost, "/order", bytes.NewReader(body))
	request.Header.Set("Authorization", "Bearer "+env.Token)
	writer := httptest.NewRecorder()

	env.Router.ServeHTTP(writer, request)

	if writer.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d, body = %s", writer.Code, http.StatusCreated, writer.Body)
	}

	var created Order.Order
	if err := json.Unmarshal(writer.Body.Bytes(), &created); err != nil {
		t.Fatal(err)
	}
	if created.ID == 0 {
		t.Fatal("order was created without an id")
	}
	if created.UserID != env.User.ID {
		t.Errorf("user id = %d, want %d", created.UserID, env.User.ID)
	}
	if len(created.Products) != len(productIds) {
		t.Errorf("products in response = %d, want %d", len(created.Products), len(productIds))
	}

	stored, err := Order.NewOrderRepository(env.Db).GetById(created.ID)
	if err != nil {
		t.Fatal(err)
	}
	if stored.UserID != env.User.ID {
		t.Errorf("stored user id = %d, want %d", stored.UserID, env.User.ID)
	}
	if len(stored.Products) != len(productIds) {
		t.Errorf("stored products = %d, want %d", len(stored.Products), len(productIds))
	}
}
