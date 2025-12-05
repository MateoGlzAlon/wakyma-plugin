package entities

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
