package controllers

import (
	"go-rest-api/models"
	"go-rest-api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func FindCustomers(c *gin.Context) {
	var customers []models.Customer
	models.DB.Find(&customers)
	utils.RespondJSON(c, http.StatusOK, customers)
}
