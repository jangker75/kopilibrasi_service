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
	// var data []models.MenuCategoryWithMenuDetails
	// models.DB.Find(&data)

	// listcategoryId := make([]uint, len(data))
	// for i, category := range data {
	// 	listcategoryId[i] = category.ID
	// }

	// var menuDetails []models.MenuDetail
	// models.DB.Where("category_id IN ?", listcategoryId).Find(&menuDetails)
	// for i := range data {
	// 	data[i].MenuDetails = []models.MenuDetail{}
	// }
	var result []models.MenuCategoryWithMenuDetails
	models.DB.Find(&result)
	// err := models.DB.Table("menu_category").
	// 	Select("menu_category.id, menu_category.created_at, menu_category.updated_at, menu_category.deleted_at, menu_category.title, menu_category.category, menu_category.description").
	// 	Joins("LEFT JOIN menu_details ON menu_category.id = menu_details.category_id").
	// 	Scan(&result).Error

	// if err != nil {
	// 	utils.RespondJSON(c, http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }

	for i := range result {
		var details []models.MenuDetail
		err := models.DB.Where("category_id = ?", result[i].ID).Find(&details).Error
		if err != nil {
			utils.RespondJSON(c, http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		result[i].MenuDetails = details
	}
	utils.RespondJSON(c, http.StatusOK, result)
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
