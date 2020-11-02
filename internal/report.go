package internal

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Report []ReportTransaction

type ReportTransaction struct {
	ID           string
	Date         time.Time
	Type         string
	Total        float64
	Quantity     int
	Commission   float64
	CurrentPrice float64
	Price        float64
	Status       string
}

func (r Report) ToCSV() string {
	result := fmt.Sprintf("%s;%s;%s;%s;%s;%s;%s;%s;%s\n", "id", "date", "type", "total", "price", "qty", "commision", "current", "status")
	for _, transaction := range r {
		result += fmt.Sprintf("%s;%s;%s;%s;%s;%d;%s;%s;%s\n",
			transaction.ID,
			dateConvert(transaction.Date),
			transaction.Type,
			dotToComma(transaction.Total),
			dotToComma(transaction.Price),
			transaction.Quantity,
			dotToComma(transaction.Commission),
			dotToComma(transaction.CurrentPrice),
			transaction.Status,
		)
	}
	return result
}

type tableRow struct {
	ID         string
	Total      float64
	Qty        int
	Commission float64
	Current    float64
	AvgPrice   float64
	Price      float64
	Revenue    float64
	Percent    float64
}

func (t tableRow) ToStrings() (result []string) {
	result = append(result, t.ID)
	result = append(result, dotToComma(t.Total))
	result = append(result, strconv.Itoa(t.Qty))
	result = append(result, dotToComma(t.Commission))
	result = append(result, dotToComma(t.Current))
	result = append(result, dotToComma(t.AvgPrice))
	result = append(result, dotToComma(t.Price))
	result = append(result, dotToComma(t.Revenue))
	result = append(result, dotToComma(t.Percent))
	return
}

func (r Report) ToTermTable(currentOnly bool) (table [][]string) {
	headers := []string{"ID", "Invested", "Qty", "Commission", "Current", "Avg Price", "Price", "Revenue", "%"}
	table = append(table, headers)
	pivotTransactions := PivotTableTransactions(&r)
	for id, pivot := range pivotTransactions {
		if currentOnly && (pivot.Quantity < 1) {
			continue
		}
		avgPrice := pivot.Total / float64(pivot.Quantity)
		price := float64(pivot.Quantity) * pivot.CurrentPrice
		revenue := pivot.Total + price
		trow := tableRow{
			ID:         id,
			Total:      pivot.Total,
			Qty:        pivot.Quantity,
			Commission: pivot.Commission,
			Current:    pivot.CurrentPrice,
			AvgPrice:   avgPrice,
			Price:      price,
			Revenue:    revenue,
			Percent:    (revenue * 100) / pivot.Total,
		}
		table = append(table, trow.ToStrings())
	}
	return
}

func dateConvert(t time.Time) string {
	return fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d", t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
}

func dotToComma(f float64) string {
	return strings.ReplaceAll(fmt.Sprintf("%.2f", f), ".", ",")
}

func PivotTableTransactions(report *Report) map[string]ReportTransaction {
	result := make(map[string]ReportTransaction)
	for _, tr := range *report {
		if tr.Status != "Done" {
			continue
		}
		if _, ok := result[tr.ID]; !ok {
			result[tr.ID] = tr
			continue
		}

		accTransaction := result[tr.ID]
		accTransaction.Total += tr.Total
		accTransaction.Quantity += tr.Quantity
		accTransaction.Commission += tr.Commission
		accTransaction.Date = time.Time{}
		accTransaction.Status = ""
		accTransaction.Type = ""
		accTransaction.Price = 0
		result[tr.ID] = accTransaction
	}
	return result
}
