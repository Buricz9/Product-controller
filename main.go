package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	dsn := "admin:admin@tcp(localhost:3306)/productdb"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Błąd połączenia:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Błąd pingowania bazy:", err)
	}

	fmt.Println("Połączono z bazą danych MariaDB!")
}
