package Database

import (
	"assaultrifle/Form"
	"assaultrifle/Utils"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	dsn := Utils.GetConfig("pg", "dsn") 
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Veritabanına bağlanılamadı: " + err.Error())
	}

	DB.AutoMigrate(&Form.User{},&Form.PriceListing{},&Form.Product{},&Form.Seller{})
	fmt.Println("✅ GORM ile MySQL bağlantısı kuruldu ve tablolar migrat edildi.")
}

