package models

import "github.com/jinzhu/gorm"

type MenuDetail struct {
	gorm.Model
	ID          uint    `json:"id" gorm:"unique;primaryKey;autoIncrement"`
	CategoryId  uint    `json:"category_id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	IsRecommend bool    `json:"isrecommend" gorm:"column:isrecommend"`
}
