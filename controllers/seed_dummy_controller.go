package controllers

import (
	"go-rest-api/models"
	"go-rest-api/utils"
	"net/http"
	"strconv"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gin-gonic/gin"
)

func SeedDummyData(c *gin.Context) {
	// Parse counter query parameter
	count, err := strconv.Atoi(c.DefaultQuery("counter", "1"))
	if err != nil || count < 1 {
		utils.RespondJSON(c, http.StatusBadRequest, gin.H{"message": "Dummy data seeded failed, invalid counter value"})
		return
	}

	// Seed Menu Details
	menuDetails := make([]models.MenuDetail, count) // Use a slice of values
	for i := range menuDetails {
		menuDetails[i] = models.MenuDetail{ // Assign values directly
			CategoryId:  uint(gofakeit.Number(1, 4)),
			Name:        gofakeit.Name(),
			IsRecommend: false,
			Price:       gofakeit.Price(1000, 30000),
		}
	}
	models.DB.Create(menuDetails) // Pass the slice of values
	utils.RespondJSON(c, http.StatusOK, gin.H{"message": "Dummy data seeded successfully", "count": count})
}
