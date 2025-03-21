package main

import (
	"log"
	"net/http"
	"product-controller/config"
	"product-controller/controller"
	"product-controller/models"
	"product-controller/repository"
	"product-controller/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	config.InitDB()

	// Migracje
	err := config.DB.AutoMigrate(
		&models.Product{},
		&models.ProductHistory{},
		&models.BlacklistWord{},
	)
	if err != nil {
		log.Fatal("Błąd migracji:", err)
	}

	// Inicjalizacja warstw
	productRepo := repository.NewProductRepository()
	blacklistRepo := repository.NewBlacklistRepository()

	productService := service.NewProductService(productRepo, blacklistRepo)
	productController := controller.NewProductController(productService)
	blacklistController := controller.NewBlacklistController(blacklistRepo)

	// Router
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// Endpointy
	r.Get("/products", productController.GetAllProducts)
	r.Post("/products", productController.AddProduct)
	r.Delete("/products/{id}", productController.DeleteProduct)
	r.Put("/products/{id}", productController.UpdateProduct)
	r.Get("/products/{id}", productController.GetProductByID)

	// Endpointy dla blacklisty
	r.Get("/blacklist", blacklistController.GetAllBlacklistWords)
	r.Post("/blacklist", blacklistController.AddBlacklistWord)
	r.Delete("/blacklist/{id}", blacklistController.DeleteBlacklistWord)

	r.Get("/products/{id}/history", productController.GetProductHistory)

	log.Println("Serwer nasłuchuje na porcie :8080")
	http.ListenAndServe(":8080", r)
}
