package hitbtc

import (
	"encoding/json"
	"fmt"
	"strings"

	"go.sadegh.io/expi/types"
	"go.sadegh.io/expi/types/event"
)

func (h *HitBTC) receive(evt *event.Event) {
	switch {
	case evt.Topic == "spot_order":
		report := &ReportResponse{}
		err := json.Unmarshal(evt.Response.([]byte), &report)
		if err != nil {
			h.receiver(&event.Event{
				Topic: "error:hitbtc",
				Event: fmt.Errorf("report response error :%v", err),
			})
			return
		}
		if report.ID == 0 {
			return
		}

		h.receiver(&event.Event{Event: report})

	case strings.Contains(evt.Topic, "candles"):
		candle := &CandleResponse{}
		err := json.Unmarshal(evt.Response.([]byte), &candle)
		if err != nil {
			h.receiver(&event.Event{
				Topic: "error:hitbtc",
				Event: fmt.Errorf("candles response error :%v", err),
			})
			return
		}

		h.receiver(&event.Event{Event: candle})

	case evt.Topic == "error":
		apiError := &types.ApiErr{}
		switch evt.Response.(type) {
		case []byte:
			err := json.Unmarshal(evt.Response.([]byte), &apiError)
			if err != nil {
				h.receiver(&event.Event{
					Topic: "error:hitbtc",
					Event: fmt.Errorf("unmarshal api error :%v", err),
				})
				return
			}

			h.receiver(&event.Event{
				Topic: "error:hitbtc",
				Event: fmt.Errorf("api error :%v", error(apiError)),
			})
		case string:
			h.receiver(&event.Event{
				Topic: "error:hitbtc",
				Event: fmt.Errorf("api error :%v", evt.Response),
			})
		}
	}
}
