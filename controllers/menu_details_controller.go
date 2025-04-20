package controllers

import (
	"go-rest-api/models"
	"go-rest-api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func FindMenuDetails(c *gin.Context) {
	var data []models.MenuDetail
	// tx := models.DB.Session(&gorm.Session{Logger: newLogger})
	models.DB.Find(&data)
	utils.RespondJSON(c, http.StatusOK, data)
}
func CreateMenuDetails(c *gin.Context) {
	var input models.MenuDetail
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondJSON(c, http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	menuDetail := models.MenuDetail{CategoryId: input.CategoryId, Name: input.Name, IsRecommend: input.IsRecommend, Price: input.Price}
	models.DB.Create(&menuDetail)
	utils.RespondJSON(c, http.StatusCreated, menuDetail)
}
