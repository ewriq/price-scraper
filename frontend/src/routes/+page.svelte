<script lang="ts">
    import { onMount } from 'svelte';
  
    // API'den beklediğimiz veri yapısı
    interface Product {
      ID: number;
      Token: string;
      Name: string;
      Content: string;
      Features: string;
      // Diğer alanlar (CreatedAt, UpdatedAt, DeletedAt) da olabilir,
      // ancak şimdilik göstermeyeceksek tanımlamaya gerek yok.
    }
  
    interface ApiResponse {
      data: Product[];
      message: string;
      status: string;
    }
  
    let products: Product[] = [];
    let errorMessage: string | null = null;
    let isLoading: boolean = true;
  
    onMount(async () => {
      try {
        const response = await fetch('http://localhost:3000/api/admin/product/list');
        if (!response.ok) {
          // HTTP hatası durumunda
          const errorData = await response.json();
          throw new Error(errorData.message || 'Veri çekilirken bir hata oluştu');
        }
  
        const result: ApiResponse = await response.json();
  
        if (result.status === 'OK') {
          products = result.data;
        } else {
          errorMessage = result.message || 'API yanıtı başarısız oldu.';
        }
      } catch (error) {
        if (error instanceof Error) {
          errorMessage = error.message;
        } else {
          errorMessage = 'Bilinmeyen bir hata oluştu.';
        }
        console.error("API hatası:", error);
      } finally {
        isLoading = false;
      }
    });
  </script>
  
  <div class="product-list-container">
    <h1>Ürün Listesi</h1>
  
    {#if isLoading}
      <p>Ürünler yükleniyor...</p>
    {:else if errorMessage}
      <p class="error-message">Hata: {errorMessage}</p>
    {:else if products.length === 0}
      <p>Gösterilecek ürün bulunamadı.</p>
    {:else}
      <ul class="product-grid">
        {#each products as product (product.ID)}
          <li class="product-card">
            <h2><a href="/{product.Token}">{product.Name}</a></h2>
            <p><strong>İçerik:</strong> {product.Content}</p>
            <p><strong>Özellikler:</strong> {product.Features}</p>
            </li>
        {/each}
      </ul>
    {/if}
  </div>
  
  <style>
    .product-list-container {
      max-width: 1200px;
      margin: 40px auto;
      padding: 20px;
      background-color: #f9f9f9;
      border-radius: 8px;
      box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
      font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
      color: #333;
    }
  
    h1 {
      text-align: center;
      color: #2c3e50;
      margin-bottom: 30px;
    }
  
    .product-grid {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
      gap: 20px;
      list-style: none;
      padding: 0;
    }
  
    .product-card {
      background-color: #fff;
      border: 1px solid #ddd;
      border-radius: 8px;
      padding: 20px;
      box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
      transition: transform 0.2s ease-in-out;
    }
  
    .product-card:hover {
      transform: translateY(-5px);
    }
  
    .product-card h2 {
      color: #3498db;
      margin-top: 0;
      margin-bottom: 10px;
      font-size: 1.5em;
    }
  
    .product-card p {
      margin-bottom: 8px;
      line-height: 1.5;
    }
  
    .product-card p strong {
      color: #555;
    }
  
    .error-message {
      color: #e74c3c;
      font-weight: bold;
      text-align: center;
    }
  </style>