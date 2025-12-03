// backend/usecases/listallinvoices/service.go
package listallinvoices

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Service struct{}

type Response struct {
	Success    bool       `json:"success"`
	Data       []Invoice  `json:"data"`
	Pagination Pagination `json:"pagination"`
}

type Invoice struct {
	ID                string   `json:"id"`
	InvoiceName       string   `json:"invoiceName"`
	InvoiceNumber     int64    `json:"invoiceNumber"`
	TotalPrice        float64  `json:"totalPrice"`
	TotalIVA          float64  `json:"totalIVA"`
	TotalPriceWithTax float64  `json:"totalPriceWithTax"`
	PaidAmount        float64  `json:"paidAmount"`
	PendingAmount     float64  `json:"pendingAmount"`
	PaymentStatus     int      `json:"paymentStatus"`
	PaymentMethod     []string `json:"paymentMethod"`
	Client            Client   `json:"client"`
	Pet               Pet      `json:"pet"`
	InvoiceDate       string   `json:"invoiceDate"`
	CreatedAt         string   `json:"createdAt"`
	UpdatedAt         string   `json:"updatedAt"`
}

type Client struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
}

type Pet struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Pagination struct {
	Limit   int  `json:"limit"`
	Skip    int  `json:"skip"`
	Count   int  `json:"count"`
	Total   int  `json:"total"`
	HasMore bool `json:"hasMore"`
}

type Params struct {
	DateFrom  string `json:"dateFrom"`
	DateUntil string `json:"dateUntil"`
	ClientID  string `json:"clientId"`
	Limit     int    `json:"limit"`
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Execute() ([]Invoice, error) {
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

	return response.Data, nil
}
