package hitbtc

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"go.sadegh.io/expi/internal/cast"
	"go.sadegh.io/expi/internal/stream"
	"go.sadegh.io/expi/types"
)

const Period1Minute string = "M1"

type CandleResponse types.Candle

func (r *CandleResponse) UnmarshalJSON(data []byte) error {
	var v map[string][]struct {
		Timestamp   int         `json:"t,required"`
		Max         interface{} `json:"h,required"`
		Min         interface{} `json:"l,required"`
		Open        interface{} `json:"o,required"`
		Close       interface{} `json:"c,required"`
		Volume      interface{} `json:"v,required"`
		VolumeQuote interface{} `json:"q,required"`
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	for symbol, candles := range v {
		candle := candles[0]
		t, _ := strconv.ParseInt(strconv.Itoa(candle.Timestamp)[:10], 10, 64)

		r.Symbol = symbol
		r.Timestamp = time.Unix(t, 0)
		r.Max = cast.ToFloat64(candle.Max.(string))
		r.Min = cast.ToFloat64(candle.Min.(string))
		r.Open = cast.ToFloat64(candle.Open.(string))
		r.Close = cast.ToFloat64(candle.Close.(string))
		r.Volume = cast.ToFloat64(candle.Volume.(string))
		r.VolumeQuote = cast.ToFloat64(candle.VolumeQuote.(string))
	}

	return nil
}

func (h *HitBTC) SubscribeCandles(period string, symbols []string) error {
	if period == "" {
		period = Period1Minute
	}

	ch := fmt.Sprintf("candles/%s", period)
	period = ""

	request := &stream.Request{
		Params: struct {
			Period  string   `json:"period,omitempty"`
			Symbols []string `json:"symbols,omitempty"`
		}{
			Period:  period,
			Symbols: symbols,
		},
		Endpoint: "/public", Channel: ch, Method: "subscribe",
	}

	err := h.stream.Request(request)
	if err != nil {
		return err
	}
	go h.stream.Receive(request)

	return nil
}
