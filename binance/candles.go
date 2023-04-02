package binance

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"go.sadegh.io/expi/internal/cast"
	"go.sadegh.io/expi/internal/stream"
	"go.sadegh.io/expi/types"
)

const Period1Minute string = "1m"

type CandleResponse types.Candle

func (r *CandleResponse) UnmarshalJSON(data []byte) error {
	var v struct {
		Symbol string `json:"s,required"`
		Candle struct {
			Timestamp   int         `json:"T,required"`
			Max         interface{} `json:"h,required"`
			Min         interface{} `json:"l,required"`
			Open        interface{} `json:"o,required"`
			Close       interface{} `json:"c,required"`
			Symbol      string      `json:"s,required"`
			Volume      interface{} `json:"v,required"`
			VolumeQuote interface{} `json:"q,required"`
		} `json:"k,required"`
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	t, _ := strconv.ParseInt(strconv.Itoa(v.Candle.Timestamp)[:10], 10, 64)

	r.Symbol = v.Symbol
	r.Timestamp = time.Unix(t, 0)
	r.Max = cast.ToFloat64(v.Candle.Max.(string))
	r.Min = cast.ToFloat64(v.Candle.Min.(string))
	r.Open = cast.ToFloat64(v.Candle.Open.(string))
	r.Close = cast.ToFloat64(v.Candle.Close.(string))
	r.Volume = cast.ToFloat64(v.Candle.Volume.(string))
	r.VolumeQuote = cast.ToFloat64(v.Candle.VolumeQuote.(string))

	return nil
}

func (b *Binance) SubscribeCandles(period string, symbols []string) error {
	if period == "" {
		period = Period1Minute
	}

	var params []string
	for _, symbol := range symbols {
		params = append(params, fmt.Sprintf("%s@kline_%s", strings.ToLower(symbol), period))
	}

	request := &stream.Request{
		Method: "SUBSCRIBE", Endpoint: b.stream.RandomKey(), Params: params,
	}

	err := b.stream.Request(request)
	if err != nil {
		return err
	}
	go b.stream.Receive(request)

	return err
}
