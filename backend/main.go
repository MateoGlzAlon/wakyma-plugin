// backend/main.go
package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/MateoGlzAlon/wakyma-plugin/usecases/listallinvoices"
)

func main() {
	invoiceService := listallinvoices.NewService()

	http.HandleFunc("/invoices", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not permitted", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		invoices := invoiceService.Execute()
		if err := json.NewEncoder(w).Encode(invoices); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	log.Println("Servidor en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
