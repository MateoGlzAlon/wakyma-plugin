// backend/usecases/createreport1/createreport1_usecase.go
package createreport1

import (
	"fmt"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"

	"github.com/MateoGlzAlon/wakyma-plugin/entities"
	"github.com/MateoGlzAlon/wakyma-plugin/usecases/listallinvoices"
)

type CreateReport1Service struct{}

var (
	pendingStatus   = 1
	completedStatus = 2
)

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

	// 1. Group invoices by day
	invoicesByDay := make(map[string][]entities.Invoice)

	for _, inv := range invoices {
		dayKey := inv.InvoiceDate

		if t, err := time.Parse(time.RFC3339, inv.InvoiceDate); err == nil {
			dayKey = t.Format("2006-01-02")
		}

		invoicesByDay[dayKey] = append(invoicesByDay[dayKey], inv)
	}

	// 2. Create one Excel file per day
	for day, dayInvoices := range invoicesByDay {
		if err := createExcelForDay(day, dayInvoices); err != nil {
			return nil, err
		}
	}

	return invoices, nil
}

func createExcelForDay(day string, invoices []entities.Invoice) error {
	// 1. Create Excel
	f := excelize.NewFile()

	cSheet := "Clinica"
	tSheet := "Tienda"
	pSheet := "Pendientes cobradas"

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
		f.SetCellValue(pSheet, cell, h)
	}

	// 3. Fill rows
	rowC := 2
	rowT := 2
	rowP := 2

	for _, inv := range invoices {
		targetSheet := cSheet
		row := rowC

		if inv.PaymentStatus == pendingStatus {
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

		// Format invoice date for display (dd/mm/yyyy)
		formattedDate := inv.InvoiceDate
		if t, err := time.Parse(time.RFC3339, inv.InvoiceDate); err == nil {
			formattedDate = t.Format("02/01/2006")
		}
		f.SetCellValue(targetSheet, fmt.Sprintf("E%d", row), formattedDate)
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

	// 4. Save file using the day in the filename
	filename := fmt.Sprintf("report_%s.xlsx", day)
	if err := f.SaveAs(filename); err != nil {
		return fmt.Errorf("error saving Excel for day %s: %w", day, err)
	}

	return nil
}

func paymentStatusString(paymentStatus int) string {
	switch paymentStatus {
	case pendingStatus:
		return "Pending"
	case completedStatus:
		return "Completed"
	default:
		return "Review"
	}
}
