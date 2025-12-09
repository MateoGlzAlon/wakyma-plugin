// backend/usecases/listallinvoices/listAllInvoices_httphandler.go
package listallinvoices

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/MateoGlzAlon/wakyma-plugin/entities"
)

func ListAllInvoicesHttpHandler(endpoint string) {
	invoiceService := NewListAllInvoicesService()

	http.HandleFunc(endpoint, func(w http.ResponseWriter, r *http.Request) {
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

		fmt.Printf("listAllInvoices_Params: %+v\n", params)

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
}
