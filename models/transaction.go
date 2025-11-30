package models

import (
	"github.com/jinzhu/gorm"
)

type Transaction struct {
	gorm.Model
	TxnNumber  string     `json:"txnNumber" gorm:"column:txn_number;unique"`
	Status     string     `json:"status"`
	CreatedAt  CustomTime `json:"createdAt" gorm:"-"`
	TotalPrice float64    `json:"total" gorm:"column:total_price"`
	Customer   string     `json:"customer"`
	Items      []Item     `json:"items" gorm:"-"`
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

type SyncTransactionRequest struct {
	ClientID     string        `json:"clientId"`
	Transactions []Transaction `json:"transactions"`
}
type SyncResponse struct {
	Message      string      `json:"message"`
	ClientID     string      `json:"clientId"`
	Transactions Transaction `json:"transactions"`
	Items        []Item      `json:"items"`
}
