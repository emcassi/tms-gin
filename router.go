package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RunRoutes(r *gin.Engine, db *gorm.DB) {

	// Products
	r.GET("/products", GetAllProducts)
	r.GET("/products/:id", GetProduct)
	r.POST("/products", CreateProduct)

	// Users
	r.GET("/users", GetAllUsers)
	r.GET("/users/:id", GetUser)
	r.POST("/users", CreateUser)
	r.DELETE("/users/:id", DeleteUser)
}
