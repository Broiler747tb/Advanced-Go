package http

import (
	"fmt"
	"net/http"
	"order-api/Configs"
	"order-api/Db"
	"order-api/Product"
	"order-api/User"
	"order-api/auth"
	"order-api/middleware"
)

var ProductRepo *ProductHandler

func StartServer() {
	conf := Configs.LoadConfig()
	db := Db.NewDb(conf)

	// Ensure the tables exist before serving, so a fresh database (e.g. a new
	// hosted Postgres) works with a single `go run ./Cmd` — no separate step.
	if err := db.AutoMigrate(&Product.Product{}, &User.User{}); err != nil {
		panic(err)
	}

	router := http.NewServeMux()

	//Repositories
	product := Product.NewProductRepository(db)
	authRepo := User.NewUserRepository(db)

	//Services
	authService := auth.NewAuthService(authRepo, conf.Auth.Secret)

	//Handler
	auth.NewAuthHandler(router, auth.AuthHandlerDeps{
		Config:      conf,
		AuthService: authService,
	})

	NewProductHandler(router, ProductHandlerDeps{
		ProductDB: product,
		Secret:    conf.Auth.Secret,
	})

	ProductRepo = &ProductHandler{ProductRepository: *Product.NewProductRepository(db)}

	server := http.Server{
		Addr:    ":8081",
		Handler: middleware.Logger(router),
	}

	fmt.Println("Listening and serving on port: 8081")
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
