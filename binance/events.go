package binance

import (
	"encoding/json"
	"fmt"

	"go.sadegh.io/expi/internal/errors"
	"go.sadegh.io/expi/types/event"
)

func (b *Binance) receive(evt *event.Event) {
	switch evt.Topic {
	case "executionReport":
		report := &ReportResponse{}
		err := json.Unmarshal(evt.Response.([]byte), &report)
		if err != nil {
			b.receiver(&event.Event{Event: err})
		}

		b.receiver(&event.Event{Event: report})

	case "kline":
		candle := &CandleResponse{}
		err := json.Unmarshal(evt.Response.([]byte), &candle)
		if err != nil {
			b.receiver(&event.Event{Event: err})
		}

		b.receiver(&event.Event{Event: candle})

	case "error":
		apiError := &errors.BinanceApiErr{}
		switch evt.Response.(type) {
		case []byte:
			err := json.Unmarshal(evt.Response.([]byte), &apiError)
			if err != nil {
				b.receiver(&event.Event{Event: err})
			}

			b.receiver(&event.Event{Event: error(apiError)})
		case string:
			b.receiver(&event.Event{Event: fmt.Errorf(evt.Response.(string))})
		}
	}
}
