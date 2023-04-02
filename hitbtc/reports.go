package hitbtc

import (
	"encoding/json"
	"time"

	"go.sadegh.io/expi/internal/stream"
	"go.sadegh.io/expi/types"
)

type ReportResponse types.Report

func (r *ReportResponse) UnmarshalJSON(data []byte) error {
	var v struct {
		ID            int64     `json:"id,required"`
		Side          string    `json:"side,required"`
		Type          string    `json:"type,required"`
		Price         float64   `json:"price,string"`
		Symbol        string    `json:"symbol,required"`
		Status        string    `json:"status,required"`
		Quantity      float64   `json:"quantity,string"`
		StopPrice     float64   `json:"stop_price,omitempty"`
		UpdatedAt     time.Time `json:"created_at,required"`
		CreatedAt     time.Time `json:"updated_at,required"`
		TimeInForce   string    `json:"time_in_force,required"`
		ClientOrderID string    `json:"client_order_id,required"`
	}

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	r.ID = v.ID
	r.Side = v.Side
	r.Type = v.Type
	r.Price = v.Price
	r.Symbol = v.Symbol
	r.Status = v.Status
	r.Quantity = v.Quantity
	r.StopPrice = v.StopPrice
	r.CreatedAt = v.CreatedAt
	r.UpdatedAt = v.UpdatedAt
	r.TimeInForce = v.TimeInForce
	r.ClientOrderID = v.ClientOrderID

	return nil
}

func (h *HitBTC) SubscribeReports() error {
	request := &stream.Request{
		Endpoint: "/trading", Method: "spot_subscribe", Authenticate: true,
	}

	err := h.stream.Request(request)
	if err != nil {
		return err
	}
	go h.stream.Receive(request)

	return nil
}
