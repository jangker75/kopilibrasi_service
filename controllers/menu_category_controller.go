package controllers

import (
	"fmt"
	"go-rest-api/models"
	"go-rest-api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func FindMenuCategories(c *gin.Context) {
	var data []models.MenuCategory
	models.DB.Find(&data)
	utils.RespondJSON(c, http.StatusOK, data)
}

func FindMenuCategoriesWithMenuDetails(c *gin.Context) {
	var data []models.MenuCategoryWithMenuDetails
	models.DB.Find(&data)
	for i := range data {
		data[i].MenuDetails = []models.MenuDetail{}
	}
	utils.RespondJSON(c, http.StatusOK, data)
}

func CreateMenuCategory(c *gin.Context) {
	var input models.MenuCategory
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondJSON(c, http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	data := models.MenuCategory{Title: input.Title, Description: input.Description, Category: input.Category}
	models.DB.Create(&data)
	utils.RespondJSON(c, http.StatusCreated, data)
}

func FindMenuCategory(c *gin.Context) {
	var data models.MenuCategory
	fmt.Println("ID:", c.Param("id"))
	err := models.DB.Raw("SELECT * FROM menu_category WHERE id = ?", c.Param("id")).Scan(&data).Error
	if err != nil {
		utils.RespondJSON(c, http.StatusNotFound, gin.H{"error": "Menu Category not found!"})
		return
	}
	utils.RespondJSON(c, http.StatusOK, data)
}

func UpdateMenuCategory(c *gin.Context) {
	var data models.MenuCategory
	if err := models.DB.Where("id = ?", c.Param("id")).First(&data).Error; err != nil {
		utils.RespondJSON(c, http.StatusNotFound, gin.H{"error": "Menu Category not found!"})
		return
	}

	var input models.MenuCategory
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondJSON(c, http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	models.DB.Model(&data).Updates(input)
	utils.RespondJSON(c, http.StatusOK, data)
}

func DeleteMenuCategory(c *gin.Context) {
	var data models.MenuCategory
	if err := models.DB.Where("id = ?", c.Param("id")).First(&data).Error; err != nil {
		utils.RespondJSON(c, http.StatusNotFound, gin.H{"error": "Menu Category not found!"})
		return
	}
	models.DB.Delete(&data)
	utils.RespondJSON(c, http.StatusNoContent, nil)
}
