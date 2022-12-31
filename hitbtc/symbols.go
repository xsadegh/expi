package hitbtc

import (
	"encoding/json"

	"go.sadegh.io/expi/internal/api"
	"go.sadegh.io/expi/models"
)

type SymbolsResponse models.Symbols

func (r *SymbolsResponse) UnmarshalJSON(data []byte) error {
	var v map[string]struct {
		Base      string  `json:"base_currency,required"`
		Quote     string  `json:"quote_currency,required"`
		Precision float64 `json:"quantity_increment,string"`
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	for i, symbol := range v {
		*r = append(*r, models.Symbol{
			ID:        i,
			Base:      symbol.Base,
			Quote:     symbol.Quote,
			Precision: symbol.Precision,
		})
	}

	return nil
}

func (h *HitBTC) GetSymbols() (models.Symbols, error) {
	var response SymbolsResponse
	err := h.api.Request(api.Request{
		Method: "GET", Endpoint: "/public/symbol",
	}, nil, &response)

	return models.Symbols(response), err
}
