package Cron

import (
	"assaultrifle/Database"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

func ScrapeByWebsite(website, token, id string) {
	var result string
	var err error

	switch {
	case strings.Contains(website, "trendyol.com"):
		result, err = Trendyol(website)
	default:
		fmt.Printf("❌ Bilinmeyen site: %s\n", website, id)
		return
	}

	if err != nil {
		fmt.Printf("🛑 [%s] (%s): HATA → %v\n", token, website, id, err)
	} else {
		err, _ := Database.CreatePriceListing(id, result, token, website)
		if err != nil {
			fmt.Printf("🛑 [%s] (%s): Veritabanına eklenirken hata → %v\n", token, website, id, err)
		} else {
			fmt.Printf("✅ [%s] (%s): SONUÇ → %s\n", result, id, result, token, website)
		}
	}
}

var notFoundKeywords = []string{
	"ürün bulunamadı",
	"sayfa mevcut değil",
	"aradığınız sayfa bulunamadı",
	"404 not found",
}

func ContainsNotFound(b string) bool {
	l := strings.ToLower(b)
	for _, k := range notFoundKeywords {
		if strings.Contains(l, k) {
			return true
		}
	}
	return false
}

func fetchHTML(url string) (string, error) {
	client := &http.Client{
		Timeout: 10 * time.Second, // Zaman aşımı
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("URL için istek oluşturulurken hata: %w", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36")
	req.Header.Set("Accept-Language", "tr-TR,tr;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("URL'ye istek atılırken hata: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP hatası: %s (Status Code: %d)", resp.Status, resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Yanıt gövdesi okunurken hata: %w", err)
	}
	return string(bodyBytes), nil
}
