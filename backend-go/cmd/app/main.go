package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/lauz1n/go-gateway/internal/repository"
	"github.com/lauz1n/go-gateway/internal/service"
	"github.com/lauz1n/go-gateway/internal/web/server"
	_ "github.com/lib/pq"
)

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return defaultValue
}

func main() {
	fmt.Println("Starting server...")
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	connnectionString := fmt.Sprintf("host=%s port=%s user%s password=%s dbname=%s sslmode=%s",
		getEnv("DB_HOST", "db"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASSWORD", "postgres"),
		getEnv("DB_NAME", "gateway"),
		getEnv("DB_SSLMODE", "disable"),
	)

	db, err := sql.Open("postgres", connnectionString)
	if err != nil {
		log.Fatal("Error connecting to database: ", err)
	}
	defer db.Close()

	accountRepository := repository.NewAccountRepository(db)
	accountService := service.NewAccountService(accountRepository)

	invoiceRepository := repository.NewInvoiceRepository(db)
	invoiceService := service.NewInvoiceService(invoiceRepository, accountService)

	port := getEnv("HTTP_PORT", "8080")
	server := server.NewServer(accountService, invoiceService, port)
	server.ConfigureRoutes()

	if err := server.Start(); err != nil {
		log.Fatal("Error starting server: ", err)
	}
}
