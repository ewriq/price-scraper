<script lang="ts">
    import { onMount } from 'svelte';
    import { page } from '$app/stores'; // SvelteKit'in URL parametrelerini almak için
  
    // --- Veri Modelleri ---
    // API'den beklediğimiz 'product' objesinin yapısı
    interface Product {
      CreatedAt: string;
      UpdatedAt: string;
      DeletedAt: string | null;
      ID: number;
      Token: string;
      Name: string;
      Content: string;
      Features: string;
    }
  
    // API'den beklediğimiz 'prices' array'indeki her bir objenin yapısı
    interface PriceListing {
      ID: number; // Go modelinizdeki 'ID' veya 'id' olabilir, API çıktısına göre kontrol edin
      CreatedAt: string;
      UpdatedAt: string;
      DeletedAt: string | null;
      product_id: string;
      seller_id: string;
      price: string;
      link: string;
      collected_at: string; // Tarih string formatında geliyor
    }
  
    // API'den gelen genel yanıt yapısı
    interface ApiResponse {
      data: {
        product: Product;
        prices: PriceListing[];
      };
      status: string; // "success" bekliyoruz
      message?: string; // Hata durumunda gelebilir
    }
  
    // --- Bileşen Durumu ---
    let productToken: string;       // URL'den alınacak ürün token'ı
    let productDetail: Product | null = null; // Ürün detayları
    let productPrices: PriceListing[] = [];  // Ürüne ait fiyatlar
    let errorMessage: string | null = null; // Hata mesajı
    let isLoading: boolean = true;      // Yükleme durumu
  
    // --- Bileşen Yüklendiğinde Çalışacak Fonksiyon ---
    onMount(async () => {
      // 1. URL'den ürün token'ını al
      // SvelteKit'te dinamik parametreler '$page.params' objesinden alınır.
      productToken = $page.params.token;
  
      try {
        // 2. API isteğini yap
        const response = await fetch(`http://localhost:3000/api/admin/product/${productToken}`);
  
        // 3. Yanıtın başarılı olup olmadığını kontrol et (HTTP status 2xx)
        if (!response.ok) {
          const errorData = await response.json();
          throw new Error(errorData.message || 'Ürün detayları çekilirken bir hata oluştu.');
        }
  
        // 4. Yanıtı JSON olarak ayrıştır
        const result: ApiResponse = await response.json();
  
        // 5. API yanıtının 'status' alanını kontrol et
        if (result.status === 'success') { // API'niz 'OK' yerine 'success' döndürüyor
          productDetail = result.data.product;
          productPrices = result.data.prices;
        } else {
          errorMessage = result.message || 'API yanıtı beklenenden farklı bir durum döndürdü.';
        }
      } catch (error) {
        // 6. Hataları yakala ve göster
        if (error instanceof Error) {
          errorMessage = error.message;
        } else {
          errorMessage = 'Bilinmeyen bir hata oluştu.';
        }
        console.error("API hatası:", error);
      } finally {
        // 7. Yükleme durumunu bitir
        isLoading = false;
      }
    });
  
    // Fiyat linkine tıklanıldığında yeni sekmede açılmasını sağlayan fonksiyon
    function openLinkInNewTab(url: string) {
      window.open(url, '_blank');
    }
  </script>
  
  <div class="product-detail-container">
    {#if isLoading}
      <p class="loading-message">Ürün detayları yükleniyor...</p>
    {:else if errorMessage}
      <p class="error-message">Hata: {errorMessage}</p>
    {:else if !productDetail}
      <p class="no-data-message">Ürün detayları bulunamadı.</p>
    {:else}
      <div class="product-info-card">
        <h1>{productDetail.Name}</h1>
        <p><strong>Token:</strong> {productDetail.Token}</p>
        <p><strong>İçerik:</strong> {productDetail.Content}</p>
        <p><strong>Özellikler:</strong> {productDetail.Features}</p>
        </div>
  
      <h2>Fiyat Geçmişi ({productPrices.length} adet)</h2>
      {#if productPrices.length === 0}
        <p class="no-prices-message">Bu ürün için henüz fiyat bilgisi bulunamadı.</p>
      {:else}
        <ul class="price-list">
          {#each productPrices as price (price.ID)}
            <li class="price-item">
              <p><strong>Fiyat:</strong> {price.price}</p>
              <p>
                <strong>Toplanma Zamanı:</strong>
                {#if price.collected_at}
                  {new Date(price.collected_at).toLocaleString('tr-TR', {
                    year: 'numeric', month: 'long', day: 'numeric',
                    hour: '2-digit', minute: '2-digit', second: '2-digit'
                  })}
                {:else}
                  Bilinmiyor
                {/if}
              </p>
              <button on:click={() => openLinkInNewTab(price.link)} class="link-button">
                Ürün Linkine Git
              </button>
            </li>
          {/each}
        </ul>
      {/if}
    {/if}
  
    <a href="/" class="back-link">← Ürün Listesine Geri Dön</a>
  </div>
  
  <style>
    .product-detail-container {
      max-width: 900px;
      margin: 40px auto;
      padding: 25px;
      background-color: #f9f9f9;
      border-radius: 12px;
      box-shadow: 0 6px 12px rgba(0, 0, 0, 0.1);
      font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
      color: #333;
    }
  
    h1 {
      text-align: center;
      color: #2c3e50;
      margin-bottom: 25px;
      font-size: 2.2em;
    }
  
    h2 {
      color: #34495e;
      margin-top: 40px;
      margin-bottom: 20px;
      border-bottom: 2px solid #eee;
      padding-bottom: 10px;
      font-size: 1.8em;
    }
  
    .product-info-card {
      background-color: #fff;
      border: 1px solid #e0e0e0;
      border-radius: 10px;
      padding: 25px;
      margin-bottom: 30px;
      box-shadow: 0 4px 8px rgba(0, 0, 0, 0.08);
    }
  
    .product-info-card p, .price-item p {
      margin-bottom: 10px;
      line-height: 1.6;
      font-size: 1.05em;
    }
  
    .product-info-card p strong, .price-item p strong {
      color: #555;
    }
  
    .price-list {
      list-style: none;
      padding: 0;
    }
  
    .price-item {
      background-color: #ffffff;
      border: 1px solid #e0e0e0;
      border-radius: 8px;
      padding: 18px;
      margin-bottom: 15px;
      box-shadow: 0 2px 5px rgba(0, 0, 0, 0.05);
      display: flex;
      flex-direction: column;
      gap: 8px;
    }
  
    .link-button {
      background-color: #3498db;
      color: white;
      border: none;
      padding: 10px 15px;
      border-radius: 5px;
      cursor: pointer;
      font-size: 0.95em;
      transition: background-color 0.3s ease;
      align-self: flex-start;
    }
  
    .link-button:hover {
      background-color: #2980b9;
    }
  
    .loading-message, .error-message, .no-data-message, .no-prices-message {
      text-align: center;
      font-size: 1.1em;
      padding: 20px;
      background-color: #ffe0b2; 
      border-radius: 8px;
      border: 1px solid #ffcc80;
      color: #e65100; 
      font-weight: bold;
    }
  
    .error-message {
      background-color: #ffebee; 
      border-color: #ffcdd2;
      color: #c62828; 
    }
  
    .back-link {
      display: block;
      margin-top: 40px;
      text-align: center;
      color: #3498db;
      text-decoration: none;
      font-weight: bold;
      font-size: 1.15em;
      padding: 12px 0;
      border-top: 2px solid #eee;
      transition: color 0.3s ease;
    }
  
    .back-link:hover {
      color: #2980b9;
      text-decoration: underline;
    }
  </style>  