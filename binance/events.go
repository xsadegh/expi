package binance

import (
	"encoding/json"
	"fmt"

	"go.sadegh.io/expi/types"
	"go.sadegh.io/expi/types/event"
)

func (b *Binance) receive(evt *event.Event) {
	switch evt.Topic {
	case "executionReport":
		report := &ReportResponse{}
		err := json.Unmarshal(evt.Response.([]byte), &report)
		if err != nil {
			b.receiver(&event.Event{
				Topic: "error:binance",
				Event: fmt.Errorf("report response error :%v", err),
			})
			return
		}
		if report.ID == 0 {
			return
		}

		b.receiver(&event.Event{Event: report})

	case "kline":
		candle := &CandleResponse{}
		err := json.Unmarshal(evt.Response.([]byte), &candle)
		if err != nil {
			b.receiver(&event.Event{
				Topic: "error:binance",
				Event: fmt.Errorf("kline response error :%v", err),
			})
			return
		}

		b.receiver(&event.Event{Event: candle})

	case "error":
		apiError := &types.ApiErr{}
		switch evt.Response.(type) {
		case []byte:
			err := json.Unmarshal(evt.Response.([]byte), &apiError)
			if err != nil {
				b.receiver(&event.Event{
					Topic: "error:binance",
					Event: fmt.Errorf("unmarshal api error :%v", err),
				})
				return
			}

			b.receiver(&event.Event{
				Topic: "error:binance",
				Event: fmt.Errorf("api error :%v", error(apiError)),
			})
		case string:
			b.receiver(&event.Event{
				Topic: "error:binance",
				Event: fmt.Errorf("api error :%v", evt.Response),
			})
		}
	}
}
