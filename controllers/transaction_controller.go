package controllers

import (
	"fmt"
	"go-rest-api/models"
	"go-rest-api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

func SyncTransactions(c *gin.Context) {
	var input models.SyncTransactionRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondJSON(c, http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Sukses bind
	allItems := []models.Item{}
	dbTxns := make([]models.Transaction, 0)
	for _, txn := range input.Transactions {
		dbTxns = append(dbTxns, txn)
		for _, item := range txn.Items {
			allItems = append(allItems, models.Item{
				TxnNumber:    txn.TxnNumber,
				Name:         item.Name,
				SKU:          item.SKU,
				Price:        item.Price,
				Qty:          item.Qty,
				Discount:     item.Discount,
				DiscountType: item.DiscountType,
				LineTotal:    item.LineTotal,
			})
		}
	}

	response := models.SyncResponse{
		Message:  "Data berhasil diterima",
		ClientID: input.ClientID,
	}
	// UPSERT transactions
	if err := models.DB.Clauses(
		clause.OnConflict{
			Columns: []clause.Column{{Name: "txn_number"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"status",
				"total_price",
				"customer",
			}),
		},
	).Create(&dbTxns).Error; err != nil {
		fmt.Println("Error upserting transactions:", err)
	}
	// Replace items: delete old items first then insert new ones
	for _, t := range input.Transactions {
		models.DB.Where("txn_number = ?", t.TxnNumber).Delete(&models.Item{})
	}
	if err := models.DB.Create(&allItems).Error; err != nil {
		fmt.Println("Error inserting items:", err)
	}
	utils.RespondJSON(c, http.StatusOK, response)
}
