package hitbtc

import (
	"encoding/json"
	"fmt"
	"time"

	"go.sadegh.io/expi/internal/api"
	"go.sadegh.io/expi/models"
)

type OrderParams struct {
	Side           string    `json:"side,required"`
	Type           string    `json:"type,omitempty"`
	Price          float64   `json:"price,string"`
	Symbol         string    `json:"symbol,required"`
	Quantity       float64   `json:"quantity,string"`
	PostOnly       bool      `json:"post_only,omitempty"`
	StopPrice      float64   `json:"stop_price,omitempty"`
	ExpireTime     time.Time `json:"expire_time,omitempty"`
	TimeInForce    string    `json:"time_in_force,omitempty"`
	ClientOrderID  string    `json:"client_order_id,omitempty"`
	StrictValidate bool      `json:"strict_validate,omitempty"`
}

type OrderResponse models.Report

func (r *OrderResponse) UnmarshalJSON(data []byte) error {
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

func (h *HitBTC) NewOrder(params OrderParams) (models.Report, error) {
	response := OrderResponse{}
	err := h.api.Request(api.Request{
		Method: "POST", Endpoint: "/spot/order", Authenticate: true,
	}, params, &response)

	return models.Report(response), err
}

func (h *HitBTC) CancelOrder(symbol, orderID string) (models.Report, error) {
	if orderID != "" {
		symbol = ""
	}
	response := OrderResponse{}
	err := h.api.Request(api.Request{
		Method: "DELETE", Endpoint: fmt.Sprintf("/spot/order/%s", orderID), Authenticate: true,
	}, struct {
		Symbol string `json:"symbol,omitempty"`
	}{symbol}, &response)

	return models.Report(response), err
}
