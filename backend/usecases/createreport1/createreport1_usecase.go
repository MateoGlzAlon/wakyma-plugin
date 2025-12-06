// backend/usecases/createreport1/createreport1_usecase.go
package createreport1

import (
	"fmt"

	"github.com/xuri/excelize/v2"

	"github.com/MateoGlzAlon/wakyma-plugin/entities"
	"github.com/MateoGlzAlon/wakyma-plugin/usecases/listallinvoices"
)

type CreateReport1Service struct{}

func NewCreateReport1Service() *CreateReport1Service {
	return &CreateReport1Service{}
}

func (s *CreateReport1Service) Execute(params entities.Params) ([]entities.Invoice, error) {

	listAllInvoicesService := listallinvoices.NewListAllInvoicesService()
	responseInvoices, err := listAllInvoicesService.Execute(params)
	if err != nil {
		return nil, err
	}

	invoices := responseInvoices.Data

	// 1. Create Excel
	f := excelize.NewFile()
	sheet := "Report"
	f.SetSheetName("Sheet1", sheet)

	// 2. Write headers
	headers := []string{
		"Numero factura",
		"Nombre cliente",
		"Nombre mascota",
		"Precio",
		"Fecha factura",
	}

	for col, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(col+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	// 3. Fill rows
	for row, inv := range invoices {
		r := row + 2 // Excel starts at 1 (and 1 more for headers)

		f.SetCellValue(sheet, fmt.Sprintf("A%d", r), inv.InvoiceName)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", r), inv.Client.Name+" "+inv.Client.Surname)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", r), inv.Pet.Name)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", r), inv.TotalPriceWithTax)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", r), inv.InvoiceDate)
	}

	// 4. Save file
	if err := f.SaveAs("report1.xlsx"); err != nil {
		return nil, fmt.Errorf("error guardando Excel: %w", err)
	}

	return invoices, nil
}
