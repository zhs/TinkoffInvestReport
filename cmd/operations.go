package main

import (
	"context"
	"fmt"
	sdk "github.com/TinkoffCreditSystems/invest-openapi-go-sdk"
	"os"
	"strings"

	"time"
)

var figiHash = make(map[string]string)
var figiOperationsHash = make(map[string]bool)

type Client struct {
	client *sdk.RestClient
}

func NewClient(token string) *Client {
	client := sdk.NewRestClient(token)
	return &Client{client: client}
}

func (c Client) SaveToFile(filename string, days int) error {
	result := fmt.Sprintf("%s;%s;%s;%s;%s;%s;%s;%s;%s\n", "id", "date", "type", "total", "price", "qty", "commision", "current", "status")

	operations, err := c.GetAllOperations(days)
	if err != nil {
		return err
	}

	for _, operation := range operations {
		if operation.FIGI == "" {
			continue
		}
		tickerOperations, err := c.OperationsByFIGI(operation.FIGI, days)
		if err != nil {
			return err
		}
		if len(tickerOperations) == 0 {
			continue
		}

		currentPriceFIGI, _ := c.GetCurrentPrice(operation.FIGI)

		for _, op := range tickerOperations {
			qty := op.Quantity
			if op.OperationType == "Sell" {
				qty = 0 - qty
			}
			result += fmt.Sprintf("%s;%s;%s;%s;%s;%d;%s;%s;%s\n", c.NameByFIGI(operation.FIGI),
				dateConvert(op.DateTime),
				op.OperationType,
				dotToComma(op.Payment),
				dotToComma(op.Price),
				qty,
				dotToComma(op.Commission.Value),
				dotToComma(currentPriceFIGI),
				op.Status,
			)
		}
	}

	err = saveToFile(filename, result)
	if err != nil {
		return err
	}

	return nil
}

func (c Client) GetAllPortfolioPositions() ([]sdk.PositionBalance, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	positions, err := c.client.PositionsPortfolio(ctx, sdk.DefaultAccount)
	if err != nil {
		return nil, err
	}
	return positions, err
}

func (c Client) OperationsByFIGI(figi string, days int) ([]sdk.Operation, error) {
	if _, ok := figiOperationsHash[figi]; ok {
		return []sdk.Operation{}, nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	operations, err := c.client.Operations(ctx, sdk.DefaultAccount, time.Now().AddDate(0, 0, 0-days), time.Now(), figi)
	if err != nil {
		return nil, err
	}
	figiOperationsHash[figi] = true
	return operations, nil
}

func (c Client) GetCurrentPrice(figi string) (currentPrice float64, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	candles, err := c.client.Candles(ctx, time.Now().Add(-time.Hour), time.Now(), sdk.CandleInterval10Min, figi)
	if err != nil {
		return
	}
	if len(candles) > 0 {
		return candles[len(candles)-1].ClosePrice, nil
	}
	return
}

func (c Client) GetAllOperations(days int) ([]sdk.Operation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	orders, err := c.client.Operations(ctx, sdk.DefaultAccount, time.Now().AddDate(0, 0, 0-days), time.Now(), "")
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (c Client) NameByFIGI(figi string) string {
	if figi == "" {
		return "[empty]"
	}
	if name, ok := figiHash[figi]; ok {
		return name
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	instrument, err := c.client.InstrumentByFIGI(ctx, figi)
	if err != nil {
		fmt.Printf("%v\n", err)
		return ""
	}

	figiHash[figi] = instrument.Name

	return instrument.Name
}

func dateConvert(t time.Time) string {
	return fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d", t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
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

func dotToComma(f float64) string {
	return strings.ReplaceAll(fmt.Sprintf("%.2f", f), ".", ",")
}
