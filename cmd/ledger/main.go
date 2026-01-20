package main

import (
	"log"
	"microservices/internal/ledger"
	"microservices/pkg/database"
	"microservices/pkg/jsonutil"
	"net/http"
	"os"
)

func main() {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = "postgres://user:password@127.0.0.1:5435/ledger?sslmode=disable"
	}

	db, err := database.Connect(dsn)
	if err != nil {
		log.Printf("Warning: Database connection failed: %v", err)
	} else {
		defer db.Close()
		log.Println("Database connection established")

		// Run migration explicitly
		schema, err := os.ReadFile("internal/ledger/schema.sql")
		if err != nil {
			log.Printf("Failed to read schema file: %v", err)
		} else {
			if _, err := db.Exec(string(schema)); err != nil {
				log.Printf("Failed to run migration: %v", err)
			} else {
				log.Println("Schema migration executed successfully")
			}
		}
	}

	repo := ledger.NewRepository(db)
	handler := &LedgerHandler{repo: repo}

	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		jsonutil.WriteJSON(w, http.StatusOK, map[string]string{
			"status":  "active",
			"service": "ledger",
		})
	})

	mux.HandleFunc("/accounts", handler.CreateAccount)

	// Simple routing for /accounts/{id}
	mux.HandleFunc("/accounts/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handler.GetAccount(w, r)
			return
		}
		jsonutil.WriteErrorJSON(w, "Not Found")
	})

	mux.HandleFunc("/transactions", handler.RecordTransaction)

	log.Println("Ledger service starting on :8083")
	if err := http.ListenAndServe(":8083", mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
