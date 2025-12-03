// backend/usecases/listallinvoices/service.go
package listallinvoices

type Invoice struct {
	ID     int     `json:"id"`
	Number string  `json:"number"`
	Total  float64 `json:"total"`
}

type Service struct {
	invoices []Invoice
}

func NewService() *Service {
	return &Service{
		invoices: []Invoice{
			{ID: 1, Number: "F-001", Total: 100.50},
			{ID: 2, Number: "F-002", Total: 250.00},
		},
	}
}

func (s *Service) Execute() []Invoice {
	return s.invoices
}
