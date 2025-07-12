package Cron

import (
	"assaultrifle/Database"
	"fmt"
	"sync"
	"time"
)


type ScrapeJob struct {
	Website   string
	Token     string
	ProductID string
}


func scrapeWorker(id int, jobs <-chan ScrapeJob, wg *sync.WaitGroup) {
	defer wg.Done() 

	for job := range jobs {
		fmt.Printf("Worker %d: İşleniyor - Website: %s, Token: %s\n", id, job.Website, job.Token)
		ScrapeByWebsite(job.Website, job.Token, job.ProductID)
		fmt.Printf("Worker %d: Tamamlandı - Website: %s\n", id, job.Website)
	}
}

func StartScrapingCron() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop() 
	const maxWorkers = 5 
	jobs := make(chan ScrapeJob, maxWorkers)
	var wg sync.WaitGroup

	for i := 1; i <= maxWorkers; i++ {
		wg.Add(1)
		go scrapeWorker(i, jobs, &wg)
	}

	for {
		<-ticker.C 
		fmt.Println("Cron tetiklendi: satıcılar çekiliyor")

		websites, tokens, productIDs, err := Database.GetSellerWebsitesAndTokens()
		if err != nil {
			fmt.Println("Hata: Satıcılar alınamadı:", err)
			continue
		}

		if len(websites) == 0 {
			fmt.Println("Bilgi: satıcı bulunamadı.")
			continue
		}


		for i, site := range websites {
			jobs <- ScrapeJob{
				Website:   site,
				Token:     tokens[i],
				ProductID: productIDs[i],
			}
		}

	}
}

