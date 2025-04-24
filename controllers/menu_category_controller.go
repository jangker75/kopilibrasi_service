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
	var result []models.MenuCategoryWithMenuDetails
	models.DB.Find(&result)
	// var categoryIdList []int
	for i := range result {
		var details []models.MenuDetail
		err := models.DB.Where("category_id = ?", result[i].ID).Find(&details).Error
		if err != nil {
			utils.RespondJSON(c, http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		result[i].MenuDetails = details

		// categoryIdList = append(categoryIdList, int(result[i].ID))
	}
	// var details []models.MenuDetail
	// models.DB.Where("category_id IN (?)", categoryIdList).Find(&details)
	// fmt.Println("Category IDs:", categoryIdList)
	// fmt.Println("Menu Details length:", len(details))
	// var menuDetailsMap = make(map[uint][]models.MenuDetail)
	// for i := 0; i < len(details); i++ {
	// 	menuDetailsMap[details[i].CategoryId] = append(menuDetailsMap[details[i].CategoryId], details[i])
	// }
	// for i := 0; i < len(result); i++ {
	// 	result[i].MenuDetails = menuDetailsMap[result[i].ID]
	// }
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
