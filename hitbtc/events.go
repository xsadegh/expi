package hitbtc

import (
	"encoding/json"
	"fmt"
	"strings"

	"go.sadegh.io/expi/internal/errors"
	"go.sadegh.io/expi/types/event"
)

func (h *HitBTC) receive(evt *event.Event) {
	switch {
	case evt.Topic == "spot_order":
		report := &ReportResponse{}
		err := json.Unmarshal(evt.Response.([]byte), &report)
		if err != nil {
			h.receiver(&event.Event{Event: err})
		}
		if report.ID == 0 {
			return
		}

		h.receiver(&event.Event{Event: report})

	case strings.Contains(evt.Topic, "candles"):
		candle := &CandleResponse{}
		err := json.Unmarshal(evt.Response.([]byte), &candle)
		if err != nil {
			h.receiver(&event.Event{Event: err})
		}

		h.receiver(&event.Event{Event: candle})

	case evt.Topic == "error":
		apiError := &errors.HitBTCApiErr{}
		switch evt.Response.(type) {
		case []byte:
			err := json.Unmarshal(evt.Response.([]byte), &apiError)
			if err != nil {
				h.receiver(&event.Event{Event: err})
			}

			h.receiver(&event.Event{Event: error(apiError)})
		case string:
			h.receiver(&event.Event{Event: fmt.Errorf(evt.Response.(string))})
		}
	}
}
