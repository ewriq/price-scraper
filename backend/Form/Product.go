package Form

import (
	"gorm.io/gorm"
)


type Product struct {
	gorm.Model
	ID    uint   `gorm:"primaryKey"`
	Token string `gorm:"unique"`
	Name string `gorm:"unique"`
	Content string `gorm:"unique"`
	Features string `gorm:"unique"`
}


type ProductBody struct {
	Name     string `json:"name"`
	Content  string `json:"content"`
	Features string `json:"features"`
	User     string `json:"user"`
}