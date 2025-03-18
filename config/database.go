package config

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DB - Instancja bazy danych
var DB *gorm.DB

// InitDB - Inicjalizacja bazy
func InitDB() {
	dsn := "admin:admin@tcp(localhost:3306)/productdb?charset=utf8mb4&parseTime=True&loc=Local"

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Błąd połączenia z bazą:", err)
	}

	log.Println("Połączono z bazą danych!")
}
