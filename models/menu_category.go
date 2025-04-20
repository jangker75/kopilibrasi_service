package models

import "github.com/jinzhu/gorm"

type MenuCategory struct {
	gorm.Model
	// Id          uint   `json:"id" gorm:"unique;primaryKey;autoIncrement"`
	Title       string `json:"title"`
	Category    string `json:"category"`
	Description string `json:"description"`
}

type MenuCategoryWithMenuDetails struct {
	MenuCategory
	MenuDetails []MenuDetail `json:"menu_details"`
}
type Tabler interface {
	TableName() string
}

func (MenuCategory) TableName() string {
	return "menu_category"
}
