// backend/main.go
package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/MateoGlzAlon/wakyma-plugin/usecases/listallinvoices"
	"github.com/joho/godotenv"
)

func main() {
	invoiceService := listallinvoices.NewService()

	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system env")
	}

	http.HandleFunc("/invoices", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not permitted", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		invoices, err := invoiceService.Execute()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := json.NewEncoder(w).Encode(invoices); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	log.Println("Server running in http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
