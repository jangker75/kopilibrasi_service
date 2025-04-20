package models

type MenuDetail struct {
	// gorm.Model
	Id          uint    `json:"id" gorm:"unique;primaryKey;autoIncrement"`
	CategoryId  uint    `json:"category_id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	IsRecommend bool    `json:"isrecommend" gorm:"column:isrecommend"`
}
