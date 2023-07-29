package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Code  string `gorm:"not null;unique" json:"code"`
	Price uint   `gorm:"not null" json:"price"`
}

func main() {

	err := DB.AutoMigrate(&Product{}, &User{})

	if err != nil {
		panic("failed to migrate database: " + err.Error())
	}

	r := gin.Default()
	RunRoutes(r, DB)
	r.Run()
}
