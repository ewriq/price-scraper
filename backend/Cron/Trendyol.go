package Cron

import (
	"encoding/json" // JSON verisini işlemek için
	"fmt"
	"regexp" // Düzenli ifadeler için
	"strings"
)

// Trendyol sayfası için gerekli olan notFoundKeywords

type SchemaOrgProduct struct {
	Type        string `json:"@type"`
	Offers      struct {
		Type         string `json:"@type"`
		Price        string `json:"price"`        
		PriceCurrency string `json:"priceCurrency"`
		Availability string `json:"availability"`  
	} `json:"offers"`
	Name        string `json:"name"`      
	Description string `json:"description"` 
}


type TrendyolResponseProps struct {
	Product struct {
		MerchantListing struct {
			WinnerVariant struct {
				Price struct {
					Text string `json:"text"`
				} `json:"price"`
			} `json:"winnerVariant"`
		} `json:"merchantListing"`
	} `json:"product"`
}

func Trendyol(url string) (string, error) {
	html, err := fetchHTML(url)
	if err != nil {
		return "", fmt.Errorf("Trendyol veri çekme hatası: %w", err)
	}

	if ContainsNotFound(html) {
		return "", fmt.Errorf("Trendyol: Ürün bulunamadı veya sayfa hatası")
	}

	reProductJSON := regexp.MustCompile(`<script type="application/ld\+json">\s*\{[^}]*"@type":\s*"Product",(.|\n)*?}\s*<\/script>`)
	matchProductJSON := reProductJSON.FindStringSubmatch(html)

	if len(matchProductJSON) > 0 {
		jsonBlock := matchProductJSON[0]
		jsonStr := strings.TrimPrefix(jsonBlock, `<script type="application/ld+json">`)
		jsonStr = strings.TrimSuffix(jsonStr, `</script>`)
		jsonStr = strings.TrimSpace(jsonStr)

		var productData SchemaOrgProduct
		err := json.Unmarshal([]byte(jsonStr), &productData)
		if err == nil {
			price := productData.Offers.Price
			if price != "" {
				return fmt.Sprintf("%s %s", price, productData.Offers.PriceCurrency), nil
			}
		} else {
			fmt.Printf("Uyarı: Trendyol Schema.org JSON ayrıştırma hatası: %v\n", err)
		}
	}

	reEnvoyProps := regexp.MustCompile(`window\["__envoy_flash-sales-banner__PROPS"\]=(.*?);`)
	matchEnvoyProps := reEnvoyProps.FindStringSubmatch(html)

	if len(matchEnvoyProps) > 1 {
		jsonStr := matchEnvoyProps[1]
		var props TrendyolResponseProps // *** Burası düzeltildi: TrendyolResponseProps kullanılıyor ***
		err := json.Unmarshal([]byte(jsonStr), &props)
		if err == nil {
			price := props.Product.MerchantListing.WinnerVariant.Price.Text
			if price != "" {
				return price, nil
			}
		} else {
			fmt.Printf("Uyarı: Trendyol __envoy_flash-sales-banner__PROPS JSON ayrıştırma hatası: %v\n", err)
		}
	}

	priceRegexp := regexp.MustCompile(`<div class="product-price-container">.*?<span class="prc-dsc">(.*?)<\/span>.*?<\/div>`)
	priceMatch := priceRegexp.FindStringSubmatch(html)

	if len(priceMatch) > 1 {
		return strings.TrimSpace(priceMatch[1]), nil
	}

	// Eğer hiçbir yöntemle fiyat bulunamazsa
	return "", fmt.Errorf("Trendyol: Fiyat bulunamadı.")
}

