package binance

import (
	"encoding/json"
	"strconv"
	"time"

	"go.sadegh.io/expi/internal/api"
	"go.sadegh.io/expi/internal/cast"
	"go.sadegh.io/expi/types"
)

type OrderParams struct {
	Side          string  `json:"side,omitempty"`
	Type          string  `json:"type,omitempty"`
	Price         float64 `json:"price,omitempty"`
	Symbol        string  `json:"symbol,omitempty"`
	Quantity      float64 `json:"quantity,omitempty"`
	StopPrice     float64 `json:"stopPrice,omitempty"`
	TimeInForce   string  `json:"timeInForce,omitempty"`
	ClientOrderID string  `json:"newClientOrderId,omitempty"`
}

type OrderResponse types.Report

func (r *OrderResponse) UnmarshalJSON(data []byte) error {
	var v struct {
		ID                  int64   `json:"orderId,required"`
		Side                string  `json:"side,required"`
		Type                string  `json:"type,required"`
		Price               float64 `json:"price,string"`
		Symbol              string  `json:"symbol,required"`
		Status              string  `json:"status,required"`
		Quantity            float64 `json:"origQty,string"`
		StopPrice           float64 `json:"stopPrice,string"`
		CreatedAt           int     `json:"transactTime,omitempty"`
		UpdatedAt           int     `json:"workingTime,omitempty"`
		TimeInForce         string  `json:"timeInForce,required"`
		ClientOrderID       string  `json:"clientOrderId,required"`
		OriginClientOrderId string  `json:"origClientOrderId,required"`
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
	if v.CreatedAt > 0 {
		r.CreatedAt = time.Unix(cast.ToInt64(strconv.Itoa(v.CreatedAt)[:10]), 0)
	}
	if v.UpdatedAt > 0 {
		r.UpdatedAt = time.Unix(cast.ToInt64(strconv.Itoa(v.UpdatedAt)[:10]), 0)
	}

	if r.Status == "CANCELED" {
		r.ClientOrderID = v.OriginClientOrderId
	}

	if r.ClientOrderID == "" {
		r.ClientOrderID = v.OriginClientOrderId
	}

	return nil
}

func (b *Binance) NewOrder(params OrderParams) (types.Report, error) {
	response := OrderResponse{}
	err := b.api.Request(api.Request{
		Method: "POST", Endpoint: "/order", Authenticate: true,
	}, params, &response)

	return types.Report(response), err
}

func (b *Binance) CancelOrder(symbol, orderID string) (types.Report, error) {
	response := OrderResponse{}
	err := b.api.Request(api.Request{
		Method: "DELETE", Endpoint: "/order", Authenticate: true,
	}, struct {
		Symbol  string `json:"symbol,required"`
		OrderID string `json:"origClientOrderId"`
	}{symbol, orderID}, &response)

	return types.Report(response), err
}
