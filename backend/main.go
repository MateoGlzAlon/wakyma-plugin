// backend/main.go
package main

import (
	"log"
	"net/http"

	"github.com/MateoGlzAlon/wakyma-plugin/usecases/createreport1"
	"github.com/MateoGlzAlon/wakyma-plugin/usecases/listallinvoices"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system env")
	}

	listallinvoices.ListAllInvoicesHttpHandler("/invoices")

	createreport1.CreateReport1HttpHandler("/invoices/report1")

	log.Println("Server running in http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
