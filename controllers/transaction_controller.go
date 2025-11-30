package controllers

import (
	"go-rest-api/models"
	"go-rest-api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SyncTransactions(c *gin.Context) {
	var input models.SyncTransactionRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondJSON(c, http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Sukses bind
	allItems := []models.Item{}
	for _, txn := range input.Transactions {
		allItems = append(allItems, txn.Items...)
	}
	transaction := models.Transaction{
		TxnNumber:  input.Transactions[0].TxnNumber,
		Status:     input.Transactions[0].Status,
		TotalPrice: input.Transactions[0].TotalPrice,
		Customer:   input.Transactions[0].Customer,
		Items:      []models.Item{},
	}
	response := models.SyncResponse{
		Message:      "Data berhasil diterima",
		ClientID:     input.ClientID,
		Transactions: transaction,
		Items:        allItems,
	}
	utils.RespondJSON(c, http.StatusOK, response)
}
