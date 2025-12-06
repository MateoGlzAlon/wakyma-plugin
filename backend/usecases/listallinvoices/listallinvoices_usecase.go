// backend/usecases/listallinvoices/service.go
package listallinvoices

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/MateoGlzAlon/wakyma-plugin/entities"
)

type ListAllInvoicesService struct{}

func NewListAllInvoicesService() *ListAllInvoicesService {
	return &ListAllInvoicesService{}
}

func (s *ListAllInvoicesService) Execute(params entities.Params) (entities.Response, error) {
	client := &http.Client{}

	var (
		apiUrl       = os.Getenv("API_URL_WAKYMA")
		apiKeyWakyma = os.Getenv("API_KEY_WAKYMA")
	)

	fmt.Printf("URL is: %+v\n", apiUrl)

	req, err := http.NewRequest("GET", apiUrl+"/invoices", nil)
	if err != nil {
		return entities.Response{}, err
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

	if params.ClientID != "" {
		q.Add("clientId", params.ClientID)
	}

	if params.Limit != 0 {
		q.Add("limit", strconv.Itoa(params.Limit))
	}

	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return entities.Response{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return entities.Response{}, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	var response entities.Response

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return entities.Response{}, err
	}

	return response, nil
}
