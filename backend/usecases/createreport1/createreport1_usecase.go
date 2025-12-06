// backend/usecases/createreport1/createreport1_usecase.go
package createreport1

import (
	"fmt"
	"strings"

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

	cSheet := "Clinica"
	tSheet := "Tienda"
	pSheet := "Pendientes"

	f.SetSheetName("Sheet1", cSheet)
	f.NewSheet(tSheet)
	f.NewSheet(pSheet)

	// 2. Write headers
	headers := []string{
		"Numero factura",
		"Nombre cliente",
		"Nombre mascota",
		"Precio",
		"Fecha factura",
		"Estado de pago",
	}

	for col, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(col+1, 1)
		f.SetCellValue(cSheet, cell, h)
		f.SetCellValue(tSheet, cell, h)
	}

	// 3. Fill rows
	rowC := 2
	rowT := 2
	rowP := 2

	for _, inv := range invoices {
		targetSheet := cSheet
		row := rowC

		if inv.PaymentStatus == 1 {
			fmt.Printf("Factura pendiente: %s\n", inv.InvoiceName)
			targetSheet = pSheet
			row = rowP
		} else if strings.HasPrefix(inv.InvoiceName, "T") {
			targetSheet = tSheet
			row = rowT
		}

		f.SetCellValue(targetSheet, fmt.Sprintf("A%d", row), inv.InvoiceName)
		f.SetCellValue(targetSheet, fmt.Sprintf("B%d", row), inv.Client.Name+" "+inv.Client.Surname)
		f.SetCellValue(targetSheet, fmt.Sprintf("C%d", row), inv.Pet.Name)
		f.SetCellValue(targetSheet, fmt.Sprintf("D%d", row), inv.TotalPriceWithTax)
		f.SetCellValue(targetSheet, fmt.Sprintf("E%d", row), inv.InvoiceDate)
		f.SetCellValue(targetSheet, fmt.Sprintf("F%d", row), paymentStatusString(inv.PaymentStatus))

		switch targetSheet {
		case cSheet:
			rowC++
		case tSheet:
			rowT++
		case pSheet:
			rowP++
		}
	}

	if err := f.SaveAs("report1.xlsx"); err != nil {
		return nil, fmt.Errorf("error guardando Excel: %w", err)
	}

	return invoices, nil
}

func paymentStatusString(paymentStatus int) string {
	switch paymentStatus {
	case 1:
		return "Pendiente"
	case 2:
		return "Completado"
	default:
		return "Revisar"
	}
}
