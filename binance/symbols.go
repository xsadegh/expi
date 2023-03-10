package binance

import (
	"encoding/json"

	"go.sadegh.io/expi/internal/api"
	"go.sadegh.io/expi/types"
)

type SymbolsResponse types.Symbols

func (r *SymbolsResponse) UnmarshalJSON(data []byte) error {
	var v struct {
		Symbols []struct {
			ID        string  `json:"symbol,required"`
			Base      string  `json:"baseAsset,required"`
			Quote     string  `json:"quoteAsset,required"`
			Precision float64 `json:"baseAssetPrecision,required"`
		} `json:"symbols,required"`
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	for _, symbol := range v.Symbols {
		*r = append(*r, types.Symbol{
			ID:        symbol.ID,
			Base:      symbol.Base,
			Quote:     symbol.Quote,
			Precision: symbol.Precision,
		})
	}

	return nil
}

func (b *Binance) GetSymbols() (types.Symbols, error) {
	response := SymbolsResponse{}
	err := b.api.Request(api.Request{
		Method: "GET", Endpoint: "/exchangeInfo",
	}, nil, &response)

	return types.Symbols(response), err
}
