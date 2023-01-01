package types

import "time"

type Asset struct {
	Currency string
	Lock     float64
	Free     float64
}

type Assets []Asset

type Symbol struct {
	ID        string
	Base      string
	Quote     string
	Precision float64
}

type Symbols []Symbol

type Candle struct {
	Min         float64
	Max         float64
	Open        float64
	Close       float64
	Volume      float64
	Symbol      string
	Timestamp   time.Time
	VolumeQuote float64
}

type Candles []Candle

type Report struct {
	ID            int64
	Side          string
	Type          string
	Price         float64
	Symbol        string
	Status        string
	Quantity      float64
	StopPrice     float64
	UpdatedAt     time.Time
	CreatedAt     time.Time
	TimeInForce   string
	ClientOrderID string
}

type Reports []Report

type Exchange interface {
	NewOrder()
	CancelOrder()
	GetSymbols()
	GetAssets()
	SubscribeReports() error
	SubscribeCandles() error
}
