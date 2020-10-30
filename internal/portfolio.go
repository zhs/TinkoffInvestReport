package internal

import (
	"context"
	sdk "github.com/TinkoffCreditSystems/invest-openapi-go-sdk"
	"time"
)

func NewPortfolio(client *sdk.RestClient) *Portfolio {
	return &Portfolio{
		client: client,
	}
}

type Portfolio struct {
	client             *sdk.RestClient
	figiHash           map[string]string
	figiOperationsHash map[string]bool
}

func (p Portfolio) GetReport(days int) (*Report, error) {
	p.figiHash = make(map[string]string)
	p.figiOperationsHash = make(map[string]bool)
	var report Report

	operations, err := p.Operations(days)
	if err != nil {
		return nil, err
	}

	for _, operation := range operations {
		if operation.FIGI == "" {
			continue
		}
		tickerOperations, err := p.OperationsByFIGI(operation.FIGI, days)
		if err != nil {
			return nil, err
		}
		if len(tickerOperations) == 0 {
			continue
		}

		currentPriceFIGI, _ := p.CurrentPrice(operation.FIGI)

		for _, op := range tickerOperations {
			qty := op.Quantity
			if op.OperationType == "Sell" {
				qty = 0 - qty
			}

			tr := ReportTransaction{
				ID:           p.NameByFIGI(operation.FIGI),
				Date:         op.DateTime,
				Type:         string(op.OperationType),
				Total:        op.Payment,
				Price:        op.Price,
				Quantity:     qty,
				Commission:   op.Commission.Value,
				CurrentPrice: currentPriceFIGI,
				Status:       string(op.Status),
			}
			report = append(report, tr)
		}
	}

	return &report, nil
}

func (p Portfolio) NameByFIGI(figi string) string {
	if figi == "" {
		return "[empty]"
	}
	if name, ok := p.figiHash[figi]; ok {
		return name
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	instrument, err := p.client.InstrumentByFIGI(ctx, figi)
	if err != nil {
		return ""
	}

	p.figiHash[figi] = instrument.Name
	return instrument.Name
}

func (p Portfolio) CurrentPrice(figi string) (currentPrice float64, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	candles, err := p.client.Candles(ctx, time.Now().Add(-time.Hour), time.Now(), sdk.CandleInterval10Min, figi)
	if err != nil {
		return
	}
	if len(candles) > 0 {
		return candles[len(candles)-1].ClosePrice, nil
	}
	return
}

func (p *Portfolio) OperationsByFIGI(figi string, days int) ([]sdk.Operation, error) {
	if _, ok := p.figiOperationsHash[figi]; ok {
		return []sdk.Operation{}, nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	operations, err := p.client.Operations(ctx, sdk.DefaultAccount, time.Now().AddDate(0, 0, 0-days), time.Now(), figi)
	if err != nil {
		return nil, err
	}
	p.figiOperationsHash[figi] = true
	return operations, nil
}

func (p Portfolio) Operations(days int) ([]sdk.Operation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	operations, err := p.client.Operations(ctx, sdk.DefaultAccount, time.Now().AddDate(0, 0, 0-days), time.Now(), "")
	if err != nil {
		return nil, err
	}
	return operations, nil
}
