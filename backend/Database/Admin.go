package Database

import (
	"assaultrifle/Form"
	"assaultrifle/Utils"
	"fmt"
	"sort"

	"gorm.io/gorm"
)

func ValidateAdminAccess(token string) (string, error) {
	var user Form.User
	result := DB.Where("token = ?", token).First(&user)
	if result.Error != nil {
		return "", result.Error
	}
	return user.Perm, nil
}


func CreateProduct(name, content, features, user string)  (error,string) {
	perm, ValidatErr := ValidateAdminAccess(user)
	if  perm != "admin" {
		return nil, "Admin değil"
	}

	if ValidatErr != nil {
		return ValidatErr, ""
	}
	
	token, errForToken := Utils.Token(10)
	if (token == "err") {
		return errForToken, ""
	}
	
	product := Form.Product{
		Name:    name,
		Content: content,
		Features: features,
		Token:  token,
	}

	err := DB.Create(&product).Error; 
	if err != nil {
		return err, ""
	}

	return nil, token
}

func CreatePriceListing(productID, price, sellerToken, link string) (error, uint) {
	var seller Form.Seller
	err := DB.First(&seller, "token = ?", sellerToken).Error
	if err != nil {
		return fmt.Errorf("Satıcı bulunamadı: %v", err), 0
	}

	priceListing := Form.PriceListing{
		ProductID:   productID,
		SellerID:    seller.Token,
		Price:       price,
		Link:        link,
	}

	err = DB.Create(&priceListing).Error
	if err != nil {
		return err, 0
	}

	return nil, priceListing.ID
}

func CreateSeller(productID, name, website, logo, user  string) (error, string) {
	perm, validateErr := ValidateAdminAccess(user)
	if perm != "admin" {
		return nil, "Admin değil"
	}

	if validateErr != nil {
		return validateErr, ""
	}

	token, tokenErr := Utils.Token(16)
	if token == "err" {
		return tokenErr, ""
	}

	seller := Form.Seller{
		Token:   token,
		Name:    name,
		Website: website,
		Logo:    logo,
		ProductID:   productID,
	}

	err := DB.Create(&seller).Error
	if err != nil {
		return err, ""
	}

	return nil, token
}


func GetSellerWebsitesAndTokens() ([]string, []string,  []string, error) {
	var sellers []Form.Seller
	err := DB.Select("website", "token", "product_id").Find(&sellers).Error
	if err != nil {
		return nil, nil, nil,err
	}

	var websites []string
	var tokens []string
	var productIDs []string
	for _, s := range sellers {
		websites = append(websites, s.Website)
		tokens = append(tokens, s.Token)
		productIDs = append(productIDs, s.ProductID)
	}

	return websites, tokens, productIDs, nil
}

func GetProductByToken(token string) (Form.Product, error) {
	var product Form.Product
	result := DB.Where("token = ?", token).First(&product)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return product, fmt.Errorf("Ürün bulunamadı: %s", token)
		}
		return product, fmt.Errorf("Ürün çekme hatası: %w", result.Error)
	}
	return product, nil
}


func GetAllProducts() ([]Form.Product, error) {
	var products []Form.Product
	result := DB.Find(&products)
	if result.Error != nil {
		return nil, fmt.Errorf("Tüm ürünleri çekme hatası: %w", result.Error)
	}
	return products, nil
}


func GetSellerByToken(token string) (Form.Seller, error) {
	var seller Form.Seller
	result := DB.Where("token = ?", token).First(&seller)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return seller, fmt.Errorf("Satıcı bulunamadı: %s", token)
		}
		return seller, fmt.Errorf("Satıcı çekme hatası: %w", result.Error)
	}
	return seller, nil
}

func GetAllSellers() ([]Form.Seller, error) {
	var sellers []Form.Seller
	result := DB.Find(&sellers)
	if result.Error != nil {
		return nil, fmt.Errorf("Tüm satıcıları çekme hatası: %w", result.Error)
	}
	return sellers, nil
}


func GetPriceListingsByProductID(productToken string) ([]Form.PriceListing, error) {
	var listings []Form.PriceListing
	result := DB.Where("product_id = ?", productToken).Find(&listings)
	if result.Error != nil {
		return nil, fmt.Errorf("Ürün fiyat listelemelerini çekme hatası: %w", result.Error)
	}
	return listings, nil
}


func GetPriceListingsBySellerID(sellerToken string) ([]Form.PriceListing, error) {
	var listings []Form.PriceListing
	result := DB.Where("seller_id = ?", sellerToken).Find(&listings)
	if result.Error != nil {
		return nil, fmt.Errorf("Satıcı fiyat listelemelerini çekme hatası: %w", result.Error)
	}
	return listings, nil
}


func GetLatestPriceListingForProduct(productToken string) (Form.PriceListing, error) {
	var listing Form.PriceListing
	result := DB.Where("product_id = ?", productToken).Order("collected_at DESC").First(&listing)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return listing, fmt.Errorf("Bu ürün için fiyat listelemesi bulunamadı: %s", productToken)
		}
		return listing, fmt.Errorf("En son fiyat listelemesini çekme hatası: %w", result.Error)
	}
	return listing, nil
}

func GetProductWithAllPrices(token string) (Form.ProductWithPrices, error) {
    var product Form.Product
    result := DB.Where("token = ?", token).First(&product)
    if result.Error != nil {
        if result.Error == gorm.ErrRecordNotFound {
            return Form.ProductWithPrices{}, fmt.Errorf("Ürün bulunamadı: %s", token)
        }
        return Form.ProductWithPrices{}, fmt.Errorf("Ürün çekme hatası: %w", result.Error)
    }

    var allPrices []Form.PriceListing
    priceResult := DB.Where("product_id = ?", product.Token).Order("collected_at DESC").Find(&allPrices)
    if priceResult.Error != nil {
        if priceResult.Error == gorm.ErrRecordNotFound {
            fmt.Printf("Uyarı: Ürün %s için fiyat listelemesi bulunamadı.\n", product.Token)
            return Form.ProductWithPrices{Product: product, Prices: []Form.PriceListing{}}, nil
        }
        return Form.ProductWithPrices{}, fmt.Errorf("Fiyat listelemelerini çekme hatası: %w", priceResult.Error)
    }

    uniqueRecentPricesMap := make(map[string]Form.PriceListing)
    for _, p := range allPrices {
        if _, ok := uniqueRecentPricesMap[p.Link]; !ok {
            uniqueRecentPricesMap[p.Link] = p
        }
    }

    var finalPrices []Form.PriceListing
    for _, p := range uniqueRecentPricesMap {
        finalPrices = append(finalPrices, p)
    }

    sort.Slice(finalPrices, func(i, j int) bool {
        return finalPrices[i].CollectedAt.After(finalPrices[j].CollectedAt)
    })

    return Form.ProductWithPrices{
        Product: product,
        Prices:  finalPrices,
    }, nil
}