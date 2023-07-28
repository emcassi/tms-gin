package main

import (
	"net/http"

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
	r.POST("/login", Login)
	r.PATCH("/users/update-avatar", AuthMiddleware(), SetAvatar)

	r.GET("/private", AuthMiddleware(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "You are authorized to access this route"})
	})

	r.GET("/current-user", AuthMiddleware(), func(c *gin.Context) {
		id, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not logged in"})
			return
		}

		var user User
		DB.First(&user, id)
		c.JSON(http.StatusOK, gin.H{"id": user.ID, "email": user.Email, "created": user.CreatedAt})
	})
}
