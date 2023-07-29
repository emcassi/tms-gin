package main

import (
	"net/http"

	"github.com/emcassi/gin-tms/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RunRoutes(r *gin.Engine, db *gorm.DB) {

	r.Static("/avatars", "./avatars")

	// Products
	r.GET("/products", GetAllProducts)
	r.GET("/products/:id", GetProduct)
	r.POST("/products", CreateProduct)

	// Users
	r.GET("/users", controllers.GetAllUsers)
	r.GET("/users/:id", GetUser)
	r.POST("/users", CreateUser)
	r.DELETE("/users/:id", DeleteUser)
	r.POST("/login", Login)
	r.PATCH("/users/update-avatar", AuthMiddleware(), SetAvatar)

	r.GET("/private", AuthMiddleware(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "You are authorized to access this route"})
	})

	r.GET("/current-user", AuthMiddleware(), func(c *gin.Context) {
		user, err := GetCurrentUser(c)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"id": user.ID, "email": user.Email, "avatar": user.Avatar, "created": user.CreatedAt})
	})
}
