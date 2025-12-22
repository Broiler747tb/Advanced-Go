package http

import (
	"fmt"
	"net/http"
	"order-api/Configs"
	"order-api/Db"
	"order-api/Product"
	"order-api/middleware"
)

var ProductRepo *ProductHandler

func StartServer() {
	conf := Configs.LoadConfig()
	deps := ProductHandlerDeps{}
	database := Db.NewDb(conf)
	router := http.NewServeMux()
	NewLinkHandler(router, deps)
	ProductRepo = &ProductHandler{ProductRepository: *Product.NewProductRepository(database)}

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
