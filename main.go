package main

import (
	"log"
	"product-controller/config"
	"product-controller/models"
)

func main() {
	// Inicjalizacja bazy
	config.InitDB()

	// Automatyczne migracje
	err := config.DB.AutoMigrate(
		&models.Product{},
		&models.ProductHistory{},
		&models.BlacklistWord{},
	)
	if err != nil {
		log.Fatal("Błąd migracji:", err)
	}

	log.Println("Migracje zakończone!")
}
