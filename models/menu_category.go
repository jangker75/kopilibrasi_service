package models

import "github.com/jinzhu/gorm"

type MenuCategory struct {
	gorm.Model
	ID          uint   `gorm:"unique;primaryKey;autoIncrement" json:"id"` // Exclude from JSON
	Title       string `json:"title"`
	Category    string `json:"category"`
	Description string `json:"description"`
}

type MenuCategoryWithMenuDetails struct {
	MenuCategory
	MenuDetails []MenuDetail `json:"menu_details" gorm:"-"`
}
type Tabler interface {
	TableName() string
}

func (MenuCategory) TableName() string {
	return "menu_category"
}
