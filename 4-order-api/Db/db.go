package Db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"order-api/Configs"
)

type Db struct {
	*gorm.DB
}

func NewDb(conf *Configs.Config) *Db {
	db, err := gorm.Open(postgres.Open(conf.Db.Dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return &Db{db}
}
