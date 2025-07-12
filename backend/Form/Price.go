package Form

import (
	"time"

	"gorm.io/gorm"
)


type ProductWithPrices struct {
	Product  Product        `json:"product"`
	Prices   []PriceListing `json:"prices"`
}
type PriceListing struct {
	gorm.Model
	ID          uint      `gorm:"primaryKey" json:"id"`
	ProductID   string    `gorm:"not null" json:"product_id"`
	SellerID    string    `gorm:"not null" json:"seller_id"`
	Price       string    `gorm:"not null" json:"price"`
	Link        string    `gorm:"type:text" json:"link"`
	CollectedAt time.Time `gorm:"not null" json:"collected_at"` 
}
