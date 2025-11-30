package models

import (
	"github.com/jinzhu/gorm"
)

type Transaction struct {
	gorm.Model
	TxnNumber       string     `json:"txnNumber" gorm:"column:txn_number;unique"`
	Status          string     `json:"status"`
	TransactionDate CustomTime `json:"transactionDate" gorm:"column:transaction_date"`
	TotalPrice      float64    `json:"total" gorm:"column:total_price"`
	Customer        string     `json:"customer"`
	Items           []Item     `json:"items" gorm:"-"`
}
type Item struct {
	gorm.Model
	TxnNumber    string  `json:"txnNumber" gorm:"column:txn_number"`
	Name         string  `json:"name"`
	SKU          string  `json:"sku"`
	Price        float64 `json:"price"`
	Qty          int     `json:"qty"`
	Discount     float64 `json:"discount"`
	DiscountType string  `json:"discountType" gorm:"column:discount_type"`
	LineTotal    float64 `json:"lineTotal" gorm:"column:line_total"`
}

// Custom table name
func (Item) TableName() string {
	return "transaction_items"
}

// IncomingTransaction represents the JSON shape used by clients when syncing
// transactions. We separate this from the DB model to avoid unmarshalling
// directly into gorm.Model (which has time.Time fields parsed with RFC3339).
type IncomingTransaction struct {
	TxnNumber  string     `json:"txnNumber"`
	Status     string     `json:"status"`
	CreatedAt  CustomTime `json:"createdAt"`
	TotalPrice float64    `json:"total"`
	Customer   string     `json:"customer"`
	Items      []Item     `json:"items"`
}

type SyncTransactionRequest struct {
	ClientID     string                `json:"clientId"`
	Transactions []IncomingTransaction `json:"transactions"`
}
type SyncResponse struct {
	Message  string `json:"message"`
	ClientID string `json:"clientId"`
}
