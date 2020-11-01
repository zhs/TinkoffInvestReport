package internal

import (
	"fmt"
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

func (r Report) ToTermTable() [][]string {
	return nil
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
