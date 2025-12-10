package createreport1

import (
	"encoding/json"
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

var (
	cashMethod     = "0"
	cardMethod     = "1"
	transferMethod = "2"
)

const reportsDir = "reports"

var (
	msgFetchingDay     = "📡 Obteniendo facturas del día %s... | "
	msgFetchedInvoices = "Se han obtenido %d facturas.\n"
	msgRateLimitWait   = "⏳ Límite de peticiones alcanzado. Esperando %d segundos para continuar...\n"
	msgTotalInvoices   = "📊 Total de facturas obtenidas: %d\n"
	msgReportsDir      = "📁 Los informes se guardarán en: %s\n"
	msgNoInvoicesDay   = "ℹ️ No hay facturas para %s, se omite la creación del informe.\n"
	msgSavedReport     = "✅ Informe guardado para %s en: %s\n"
)

func NewCreateReport1Service() *CreateReport1Service {
	return &CreateReport1Service{}
}

func (s *CreateReport1Service) Execute(params entities.Params) ([]entities.Invoice, error) {
	listAllInvoicesService := listallinvoices.NewListAllInvoicesService()

	const dateLayout = "02-01-2006"

	fromDate, err := time.Parse(dateLayout, params.DateFrom)
	if err != nil {
		return nil, fmt.Errorf("error parsing DateFrom (%s): %w", params.DateFrom, err)
	}

	toDate, err := time.Parse(dateLayout, params.DateUntil)
	if err != nil {
		return nil, fmt.Errorf("error parsing DateUntil (%s): %w", params.DateUntil, err)
	}

	if toDate.Before(fromDate) {
		return nil, fmt.Errorf("DateUntil (%s) is before DateFrom (%s)", params.DateUntil, params.DateFrom)
	}

	var allInvoices []entities.Invoice

	for d := fromDate; !d.After(toDate); d = d.AddDate(0, 0, 1) {
		dayStr := d.Format(dateLayout)

		dayParams := params
		dayParams.DateFrom = dayStr
		dayParams.DateUntil = dayStr

		for {
			fmt.Printf(msgFetchingDay, dayStr)

			responseInvoices, err := listAllInvoicesService.Execute(dayParams)
			if err == nil {
				fmt.Printf(msgFetchedInvoices, len(responseInvoices.Data))
				if len(responseInvoices.Data) > 0 {
					allInvoices = append(allInvoices, responseInvoices.Data...)
				}
				break
			}

			if strings.Contains(err.Error(), "RATE_LIMIT") || strings.Contains(err.Error(), "429") {
				waitSeconds := extractRetryAfter(err.Error())
				if waitSeconds == 0 {
					waitSeconds = 60
				}

				fmt.Printf(msgRateLimitWait, waitSeconds)
				time.Sleep(time.Duration(waitSeconds) * time.Second)
				continue
			}

			return nil, fmt.Errorf("error fetching invoices for %s: %w", dayStr, err)
		}
	}

	fmt.Printf(msgTotalInvoices, len(allInvoices))
	invoices := allInvoices

	if err := os.MkdirAll(reportsDir, 0o755); err != nil {
		return nil, fmt.Errorf("could not create reports directory %q: %w", reportsDir, err)
	}

	absReportsDir, err := filepath.Abs(reportsDir)
	if err != nil {
		return nil, err
	}

	fmt.Printf(msgReportsDir, absReportsDir)

	invoicesByDay := make(map[string][]entities.Invoice)

	for _, inv := range invoices {
		dayKey := inv.InvoiceDate

		if t, err := time.Parse(time.RFC3339, inv.InvoiceDate); err == nil {
			dayKey = t.Format("02-01-2006")
		}

		invoicesByDay[dayKey] = append(invoicesByDay[dayKey], inv)
	}

	for day, dayInvoices := range invoicesByDay {
		if len(dayInvoices) == 0 {
			fmt.Printf(msgNoInvoicesDay, day)
			continue
		}

		filePath, err := createExcelForDay(day, dayInvoices, reportsDir)
		if err != nil {
			return nil, err
		}

		fmt.Printf(msgSavedReport, day, filePath)
	}

	return invoices, nil
}

func extractRetryAfter(errMsg string) int {
	start := strings.Index(errMsg, "{")
	if start == -1 {
		return 0
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(errMsg[start:]), &parsed); err != nil {
		return 0
	}

	if val, ok := parsed["retryAfter"].(float64); ok {
		return int(val)
	}

	return 0
}

func createExcelForDay(day string, invoices []entities.Invoice, reportsDir string) (string, error) {
	f := excelize.NewFile()

	cSheet := "Clinica"
	tSheet := "Tienda"
	pSheet := "Pendientes de cobrar"

	if err := f.SetSheetName("Sheet1", cSheet); err != nil {
		return "", fmt.Errorf("error creating sheet: %w", err)
	}

	if _, err := f.NewSheet(tSheet); err != nil {
		return "", fmt.Errorf("error creating sheet: %w", err)
	}

	if _, err := f.NewSheet(pSheet); err != nil {
		return "", fmt.Errorf("error creating sheet: %w", err)
	}

	headers := []string{
		"Numero factura",
		"Nombre cliente",
		"Nombre mascota",
		"Precio",
		"Fecha factura",
		"Método de pago",
		"Estado de pago",
	}

	for col, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(col+1, 1)

		if err := f.SetCellValue(cSheet, cell, h); err != nil {
			return "", fmt.Errorf("error setting header value: %w", err)
		}
		if err := f.SetCellValue(tSheet, cell, h); err != nil {
			return "", fmt.Errorf("error setting header value: %w", err)
		}
		if err := f.SetCellValue(pSheet, cell, h); err != nil {
			return "", fmt.Errorf("error setting header value: %w", err)
		}
	}

	rowC, rowT, rowP := 2, 2, 2

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

		if err := f.SetCellValue(targetSheet, fmt.Sprintf("A%d", row), inv.InvoiceName); err != nil {
			return "", fmt.Errorf("error setting cell value: %w", err)
		}
		if err := f.SetCellValue(targetSheet, fmt.Sprintf("B%d", row), inv.Client.Name+" "+inv.Client.Surname); err != nil {
			return "", fmt.Errorf("error setting cell value: %w", err)
		}
		if err := f.SetCellValue(targetSheet, fmt.Sprintf("C%d", row), inv.Pet.Name); err != nil {
			return "", fmt.Errorf("error setting cell value: %w", err)
		}
		if err := f.SetCellValue(targetSheet, fmt.Sprintf("D%d", row), inv.TotalPriceWithTax); err != nil {
			return "", fmt.Errorf("error setting cell value: %w", err)
		}

		formattedDate := inv.InvoiceDate
		if t, err := time.Parse(time.RFC3339, inv.InvoiceDate); err == nil {
			formattedDate = t.Format("02/01/2006")
		}

		if err := f.SetCellValue(targetSheet, fmt.Sprintf("E%d", row), formattedDate); err != nil {
			return "", fmt.Errorf("error setting cell value: %w", err)
		}

		methodCode := ""
		if len(inv.PaymentMethod) > 0 {
			methodCode = inv.PaymentMethod[0]
		}

		if err := f.SetCellValue(targetSheet, fmt.Sprintf("F%d", row), paymentMethodString(methodCode)); err != nil {
			return "", fmt.Errorf("error setting cell value: %w", err)
		}
		if err := f.SetCellValue(targetSheet, fmt.Sprintf("G%d", row), paymentStatusString(inv.PaymentStatus)); err != nil {
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

func paymentMethodString(paymentMethod string) string {
	switch paymentMethod {
	case cashMethod:
		return "Efectivo"
	case cardMethod:
		return "Tarjeta"
	case transferMethod:
		return "Transferencia"
	default:
		return "Otro"
	}
}
