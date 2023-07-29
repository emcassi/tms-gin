package controllers

import (
	"net/http"

	"github.com/emcassi/gin-tms/global"
	"github.com/emcassi/gin-tms/models"

	"github.com/gin-gonic/gin"
)

func GetAllUsers(c *gin.Context) {
	var users []models.User
	global.DB.Find(&users)
	c.JSON(http.StatusOK, users)
}
