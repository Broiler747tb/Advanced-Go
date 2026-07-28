package main

import (
	"order-api/Order"
	"order-api/Product"
	"order-api/User"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(postgres.Open(os.Getenv("DSN")), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&Product.Product{}, &User.User{}, &Order.Order{})
	if err != nil {
		panic(err)
	}
}
