// backend/usecases/listallinvoices/createreport1_httphandler.go
package createreport1

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/MateoGlzAlon/wakyma-plugin/entities"
)

func CreateReport1HttpHandler(endpoint string) {
	reportService := NewCreateReport1Service()

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

		w.Header().Set("Content-Type", "application/json")
		invoices, err := reportService.Execute(params)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := json.NewEncoder(w).Encode(invoices); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}
