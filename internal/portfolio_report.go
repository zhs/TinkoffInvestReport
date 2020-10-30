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
	Price        float64
	Quantity     int
	Commission   float64
	CurrentPrice float64
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

func dateConvert(t time.Time) string {
	return fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d", t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
}

func dotToComma(f float64) string {
	return strings.ReplaceAll(fmt.Sprintf("%.2f", f), ".", ",")
}
