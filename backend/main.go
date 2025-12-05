// backend/main.go
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/MateoGlzAlon/wakyma-plugin/entities"
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

		// Parameters
		q := r.URL.Query()

		dateFrom := q.Get("dateFrom")
		dateUntil := q.Get("dateUntil")
		clientID := q.Get("clientId")

		var limit int
		if limitStr := q.Get("limit"); limitStr != "" {
			if l, err := strconv.Atoi(limitStr); err == nil {
				limit = l
			}
		}

		params := entities.Params{
			DateFrom:  dateFrom,
			DateUntil: dateUntil,
			ClientID:  clientID,
			Limit:     limit,
		}

		w.Header().Set("Content-Type", "application/json")
		invoices, err := invoiceService.Execute(params)
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
