package main

import (
	"base/internal"
	sdk "github.com/TinkoffCreditSystems/invest-openapi-go-sdk"
	"github.com/subosito/gotenv"

	"os"
)

func main() {
	_ = gotenv.Load(".env")
	token := os.Getenv("token")

	client := sdk.NewRestClient(token)
	portfolio := internal.NewPortfolio(client)
	report, err := portfolio.GetReport(100)
	if err != nil {
		println(err.Error())
		return
	}

	if err = saveToFile("report.csv", report.ToCSV()); err != nil {
		println(err.Error())
	}
}

func saveToFile(filename, s string) error {
	f, err := os.Create(filename)
	defer f.Close()
	if err != nil {
		return err
	}
	_, err = f.WriteString(s)
	if err != nil {
		return err
	}
	return nil
}
