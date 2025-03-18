package main

import (
	"log"
	"product-controller/config"
	"product-controller/models"
	"product-controller/repository"
	"product-controller/service"
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

	log.Println("Migracje zakończone!")

	// Inicjalizacja repozytoriów
	productRepo := repository.NewProductRepository()
	blacklistRepo := repository.NewBlacklistRepository()

	// Dodaj zabronione słowo (na test)
	err = blacklistRepo.AddBlacklistWord(&models.BlacklistWord{Word: "zakazany"})
	if err != nil {
		log.Fatal("Błąd dodawania blacklisty:", err)
	}

	// Inicjalizacja serwisu
	productService := service.NewProductService(productRepo, blacklistRepo)

	// Produkt testowy z zakazanym słowem
	newProduct := &models.Product{
		Name:        "Zakazany laptop",
		Description: "Nie powinien się dodać",
		Price:       999.99,
		Quantity:    2,
	}

	err = productService.AddProduct(newProduct)
	if err != nil {
		log.Println("Błąd dodania produktu:", err)
	} else {
		log.Println("Dodano produkt:", newProduct)
	}
}
