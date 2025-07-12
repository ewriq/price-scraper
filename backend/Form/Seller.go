package Form

import "gorm.io/gorm"

type Seller struct {
	gorm.Model
	Token   string `gorm:"type:text"`
	Name    string `gorm:"type:text;not null;unique"`
	Website string `gorm:"type:text"`
	Logo    string `gorm:"type:text"`
	ProductID string `gorm:"type:text"`
}

type SellerBody struct {
	Name    string `json:"name"`
	Website string `json:"website"`
	Logo    string `json:"logo"`
	User    string `json:"user"`
	ProductID string `gorm:"type:text"`
}
