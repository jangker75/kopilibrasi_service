package controllers

import (
	"fmt"
	"go-rest-api/models"
	"go-rest-api/utils"
	"net/http"
	"net/url"
	"time"

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
	// Use a map to deduplicate transactions by txn_number (keep last occurrence)
	// This prevents "ON CONFLICT DO UPDATE command cannot affect row a second time" error
	txnMap := make(map[string]models.Transaction)
	txnOrder := make([]string, 0) // preserve order of first occurrence
	for _, txn := range input.Transactions {
		if _, exists := txnMap[txn.TxnNumber]; !exists {
			txnOrder = append(txnOrder, txn.TxnNumber)
		}
		// map incoming transaction to DB model (last occurrence wins)
		txnMap[txn.TxnNumber] = models.Transaction{
			TxnNumber:       txn.TxnNumber,
			TransactionDate: txn.CreatedAt,
			Status:          txn.Status,
			TotalPrice:      txn.TotalPrice,
			Customer:        txn.Customer,
			PaymentMethod:   txn.PaymentMethod,
			Notes:           txn.Notes,
		}
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
	// Build deduplicated slice preserving order
	dbTxns := make([]models.Transaction, 0, len(txnMap))
	for _, txnNum := range txnOrder {
		dbTxns = append(dbTxns, txnMap[txnNum])
	}

	response := models.SyncResponse{
		Message:  "Data berhasil diterima",
		ClientID: input.ClientID,
	}
	// UPSERT transactions one by one to avoid primary key conflicts
	for _, txn := range dbTxns {
		if err := models.DB.Clauses(
			clause.OnConflict{
				Columns: []clause.Column{{Name: "txn_number"}},
				DoUpdates: clause.AssignmentColumns([]string{
					"status",
					"total_price",
					"customer",
					"payment_method",
					"transaction_date",
					"notes",
					"updated_at",
				}),
			},
		).Create(&txn).Error; err != nil {
			fmt.Println("Error upserting transaction:", txn.TxnNumber, err)
		}
	}
	// Replace items: delete old items first then insert new ones
	// optimize: delete items for all txn_numbers in one query to avoid N+1 deletes
	txnNumbersToDelete := make([]string, 0, len(input.Transactions))
	for _, t := range input.Transactions {
		txnNumbersToDelete = append(txnNumbersToDelete, t.TxnNumber)
	}
	if len(txnNumbersToDelete) > 0 {
		models.DB.Where("txn_number IN ?", txnNumbersToDelete).Delete(&models.Item{})
	}
	if err := models.DB.Create(&allItems).Error; err != nil {
		fmt.Println("Error inserting items:", err)
	}
	utils.RespondJSON(c, http.StatusOK, response)
}

// ListTransactions returns list of transactions with item details and supports filters
// Query params: datefrom (YYYY-MM-DD), dateto (YYYY-MM-DD), status, customer
func ListTransactions(c *gin.Context) {
	// parse query params
	q := c.Request.URL.Query()
	dateFrom := q.Get("datefrom")
	dateTo := q.Get("dateto")
	status := q.Get("status")
	customer := q.Get("customer")

	db := models.DB.Table("transactions").Select("id,created_at, updated_at, txn_number, status, transaction_date, customer, total_price, payment_method, notes")

	// apply filters
	if dateFrom == "" && dateTo == "" {
		// default: last 1 month to avoid very large queries
		now := time.Now()
		from := now.AddDate(0, -1, 0)
		db = db.Where("transaction_date >= ? AND transaction_date <= ?", from, now)
	} else {
		if dateFrom != "" {
			if t, err := time.Parse("2006-01-02", dateFrom); err == nil {
				db = db.Where("transaction_date >= ?", t)
			}
		}
		if dateTo != "" {
			if t, err := time.Parse("2006-01-02", dateTo); err == nil {
				// include the whole dateTo day
				t = t.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
				db = db.Where("transaction_date <= ?", t)
			}
		}
	}
	if status != "" {
		db = db.Where("status = ?", status)
	}
	if customer != "" {
		decoded, _ := url.QueryUnescape(customer)
		db = db.Where("customer = ?", decoded)
	}

	var rows []models.Transaction
	if err := db.Order("id desc").Scan(&rows).Error; err != nil {
		utils.RespondJSON(c, http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// build result using existing models.Transaction and models.Item
	// collect txn numbers
	txnNumbers := make([]string, 0, len(rows))
	for _, r := range rows {
		txnNumbers = append(txnNumbers, r.TxnNumber)
	}

	allItems := []models.Item{}
	if len(txnNumbers) > 0 {
		models.DB.Where("txn_number IN ?", txnNumbers).Find(&allItems)
	}

	// group items by txn_number
	itemsByTxn := make(map[string][]models.Item)
	for _, it := range allItems {
		itemsByTxn[it.TxnNumber] = append(itemsByTxn[it.TxnNumber], it)
	}

	// prepare response DTO to include createdAt and updatedAt using CustomTime
	type txnResp struct {
		Id              uint              `json:"id"`
		TxnNumber       string            `json:"txnNumber"`
		Status          string            `json:"status"`
		CreatedAt       models.CustomTime `json:"createdAt"`
		UpdatedAt       models.CustomTime `json:"updatedAt"`
		TransactionDate models.CustomTime `json:"transactionDate"`
		Customer        string            `json:"customer"`
		PaymentMethod   string            `json:"paymentMethod"`
		Notes           string            `json:"notes"`
		Total           float64           `json:"total"`
		Items           []models.Item     `json:"items"`
	}

	result := make([]txnResp, 0, len(rows))
	for _, r := range rows {
		resp := txnResp{
			Id:              r.ID,
			TxnNumber:       r.TxnNumber,
			Status:          r.Status,
			CreatedAt:       models.CustomTime{Time: r.CreatedAt},
			UpdatedAt:       models.CustomTime{Time: r.UpdatedAt},
			TransactionDate: r.TransactionDate,
			Customer:        r.Customer,
			PaymentMethod:   r.PaymentMethod,
			Total:           r.TotalPrice,
			Notes:           r.Notes,
			Items:           itemsByTxn[r.TxnNumber],
		}
		result = append(result, resp)
	}

	utils.RespondJSON(c, http.StatusOK, gin.H{"transactions": result})
}
