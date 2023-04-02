package binance

import (
	"encoding/json"
	"strconv"
	"time"

	"go.sadegh.io/expi/internal/cast"
	"go.sadegh.io/expi/internal/stream"
	"go.sadegh.io/expi/types"
)

type ReportResponse types.Report

func (r *ReportResponse) UnmarshalJSON(data []byte) error {
	var v struct {
		ID            int64   `json:"i,required"`
		Side          string  `json:"S,required"`
		Type          string  `json:"o,required"`
		Price         float64 `json:"p,string"`
		Symbol        string  `json:"s,required"`
		Status        string  `json:"x,required"`
		Quantity      float64 `json:"q,string"`
		StopPrice     float64 `json:"P,string"`
		CreatedAt     int     `json:"O,required"`
		UpdatedAt     int     `json:"T,required"`
		TimeInForce   string  `json:"f,required"`
		ClientOrderID string  `json:"c,required"`
	}

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	if v.UpdatedAt == -1 {
		v.UpdatedAt = v.CreatedAt
	}

	r.ID = v.ID
	r.Side = v.Side
	r.Type = v.Type
	r.Price = v.Price
	r.Symbol = v.Symbol
	r.Status = v.Status
	r.Quantity = v.Quantity
	r.StopPrice = v.StopPrice
	r.TimeInForce = v.TimeInForce
	r.ClientOrderID = v.ClientOrderID
	r.CreatedAt = time.Unix(cast.ToInt64(strconv.Itoa(v.CreatedAt)[:10]), 0)
	r.UpdatedAt = time.Unix(cast.ToInt64(strconv.Itoa(v.UpdatedAt)[:10]), 0)

	return nil
}

func (b *Binance) SubscribeReports() error {
	request := &stream.Request{Endpoint: "/" + b.stream.ListenKey()}

	err := b.stream.Request(request)
	if err != nil {
		return err
	}
	go b.stream.Receive(request)

	return err
}
