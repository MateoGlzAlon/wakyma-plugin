package createreport1

import (
	"fmt"
	"os"
	"path/filepath"
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

const reportsDir = "reports"

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

	// 1. Ensure directory exists
	if err := os.MkdirAll(reportsDir, 0o755); err != nil {
		return nil, fmt.Errorf("could not create reports directory %q: %w", reportsDir, err)
	}

	fmt.Println("Reports will be saved in:", reportsDir)

	// 2. Group invoices by day
	invoicesByDay := make(map[string][]entities.Invoice)

	for _, inv := range invoices {
		dayKey := inv.InvoiceDate

		if t, err := time.Parse(time.RFC3339, inv.InvoiceDate); err == nil {
			dayKey = t.Format("02-01-2006")
		}

		invoicesByDay[dayKey] = append(invoicesByDay[dayKey], inv)
	}

	// 3. Create one Excel file per day
	for day, dayInvoices := range invoicesByDay {
		filePath, err := createExcelForDay(day, dayInvoices, reportsDir)
		if err != nil {
			return nil, err
		}
		// 4. Print where the file has been saved
		fmt.Printf("Saved report for %s at: %s\n", day, filePath)
	}

	return invoices, nil
}

func createExcelForDay(day string, invoices []entities.Invoice, reportsDir string) (string, error) {
	// 1. Create Excel file
	f := excelize.NewFile()

	cSheet := "Clinica"
	tSheet := "Tienda"
	pSheet := "Pendientes cobradas"

	err := f.SetSheetName("Sheet1", cSheet)
	if err != nil {
		return "", fmt.Errorf("error creating sheet: %w", err)
	}
	_, err = f.NewSheet(tSheet)
	if err != nil {
		return "", fmt.Errorf("error creating sheet: %w", err)
	}
	_, err = f.NewSheet(pSheet)
	if err != nil {
		return "", fmt.Errorf("error creating sheet: %w", err)
	}

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
		err := f.SetCellValue(cSheet, cell, h)
		if err != nil {
			return "", fmt.Errorf("error setting header value: %w", err)
		}
		err = f.SetCellValue(tSheet, cell, h)
		if err != nil {
			return "", fmt.Errorf("error setting header value: %w", err)
		}
		err = f.SetCellValue(pSheet, cell, h)
		if err != nil {
			return "", fmt.Errorf("error setting header value: %w", err)
		}
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

		err := f.SetCellValue(targetSheet, fmt.Sprintf("A%d", row), inv.InvoiceName)
		if err != nil {
			return "", fmt.Errorf("error setting cell value: %w", err)
		}
		err = f.SetCellValue(targetSheet, fmt.Sprintf("B%d", row), inv.Client.Name+" "+inv.Client.Surname)
		if err != nil {
			return "", fmt.Errorf("error setting cell value: %w", err)
		}
		err = f.SetCellValue(targetSheet, fmt.Sprintf("C%d", row), inv.Pet.Name)
		if err != nil {
			return "", fmt.Errorf("error setting cell value: %w", err)
		}
		err = f.SetCellValue(targetSheet, fmt.Sprintf("D%d", row), inv.TotalPriceWithTax)
		if err != nil {
			return "", fmt.Errorf("error setting cell value: %w", err)
		}

		// Format invoice date for display (dd/mm/yyyy)
		formattedDate := inv.InvoiceDate
		if t, err := time.Parse(time.RFC3339, inv.InvoiceDate); err == nil {
			formattedDate = t.Format("02/01/2006")
		}
		err = f.SetCellValue(targetSheet, fmt.Sprintf("E%d", row), formattedDate)
		if err != nil {
			return "", fmt.Errorf("error setting cell value: %w", err)
		}
		err = f.SetCellValue(targetSheet, fmt.Sprintf("F%d", row), paymentStatusString(inv.PaymentStatus))
		if err != nil {
			return "", fmt.Errorf("error setting cell value: %w", err)
		}

		switch targetSheet {
		case cSheet:
			rowC++
		case tSheet:
			rowT++
		case pSheet:
			rowP++
		}
	}

	filename := fmt.Sprintf("report1_%s.xlsx", day)
	filePath := filepath.Join(reportsDir, filename)

	if err := f.SaveAs(filePath); err != nil {
		return "", fmt.Errorf("error saving Excel for day %s: %w", day, err)
	}

	return filePath, nil
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
