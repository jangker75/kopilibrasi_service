package models

type MenuCategory struct {
	ID          uint   `json:"id" gorm:"primary_key"`
	Title       string `json:"title"`
	Category    string `json:"category"`
	Description string `json:"description"`
}
type Tabler interface {
	TableName() string
}

func (MenuCategory) TableName() string {
	return "menu_category"
}
