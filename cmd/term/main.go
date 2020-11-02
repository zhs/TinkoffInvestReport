package main

import (
	"base/internal"
	sdk "github.com/TinkoffCreditSystems/invest-openapi-go-sdk"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/subosito/gotenv"
	"log"
	"os"
)

func main() {
	_ = gotenv.Load(".env")
	token := os.Getenv("token")

	client := sdk.NewRestClient(token)
	portfolio := internal.NewPortfolio(client)
	report, err := portfolio.GetReport(360)
	if err != nil {
		println(err.Error())
		return
	}

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	table1 := widgets.NewTable()
	table1.Rows = report.ToTermTable(true)
	table1.TextStyle = ui.NewStyle(ui.ColorWhite)
	table1.SetRect(0, 0, 130, 200)

	table1.RowStyles[0] = ui.NewStyle(ui.ColorWhite, ui.ColorBlack, ui.ModifierBold)
	table1.FillRow = true
	table1.RowSeparator = false

	ui.Render(table1)
	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		}
	}
}
