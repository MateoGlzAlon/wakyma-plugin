// backend/usecases/listallinvoices/service.go
package listallinvoices

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	. "github.com/MateoGlzAlon/wakyma-plugin/entities"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Execute(params Params) ([]Invoice, error) {
	client := &http.Client{}

	var (
		apiUrl       = os.Getenv("API_URL_WAKYMA")
		apiKeyWakyma = os.Getenv("API_KEY_WAKYMA")
	)

	fmt.Printf("URL is: %+v\n", apiUrl)

	req, err := http.NewRequest("GET", apiUrl+"/invoices", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKeyWakyma))
	req.Header.Set("Content-Type", "application/json")

	q := req.URL.Query()

	if params.DateFrom != "" {
		q.Add("dateFrom", params.DateFrom)
	}

	if params.DateUntil != "" {
		q.Add("dateUntil", params.DateUntil)
	}

	if params.DateUntil != "" {
		q.Add("clientId", params.ClientID)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	var response Response

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	fmt.Printf("Response is: %+v\n", response)

	filteredInvoices := filterInvoices(response.Data, Params{})

	return filteredInvoices, nil
}

func filterInvoices(invoices []Invoice, params Params) []Invoice {
	filtered := []Invoice{}
	for _, invoice := range invoices {
		if params.DateFrom != "" && invoice.InvoiceDate < params.DateFrom {
			fmt.Printf("Filtering invoice: %+v\n because its date is earlier than dateFrom\n", invoice)
			continue
		}
		if params.DateUntil != "" && invoice.InvoiceDate > params.DateUntil {
			fmt.Printf("Filtering invoice: %+v\n because its date is later than dateUntil\n", invoice)
			continue
		}
		filtered = append(filtered, invoice)
	}
	return filtered
}
