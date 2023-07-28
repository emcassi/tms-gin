package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetAllProducts(c *gin.Context) {
	var products []Product
	DB.Find(&products)
	c.JSON(http.StatusOK, products)
}

func GetProduct(c *gin.Context) {
	var product Product
	DB.First(&product, c.Param("id"))
	c.JSON(http.StatusOK, product)
}

func CreateProduct(c *gin.Context) {
	var product Product
	err := c.BindJSON(&product)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check for the presence of 'code' and 'price' fields
	if product.Code == "" || product.Price == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Code and Price fields are required"})
		return
	}

	if !isCodeUnique(DB, product.Code) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Code must be unique"})
		return
	}

	DB.Create(&product)
	c.JSON(http.StatusOK, product)
}

// isCodeUnique checks if a given code is already present in the database.
func isCodeUnique(db *gorm.DB, code string) bool {
	var count int64
	db.Model(&Product{}).Where("code = ?", code).Count(&count)
	return count == 0
}
