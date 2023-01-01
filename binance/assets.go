package binance

import (
	"encoding/json"

	"go.sadegh.io/expi/internal/api"
	"go.sadegh.io/expi/types"
)

type AssetsResponse types.Assets

func (r *AssetsResponse) UnmarshalJSON(data []byte) error {
	var v struct {
		Assets []struct {
			Free     float64 `json:"free,string"`
			Lock     float64 `json:"locked,string"`
			Currency string  `json:"asset,required"`
		} `json:"balances,required"`
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	for _, asset := range v.Assets {
		*r = append(*r, types.Asset{
			Currency: asset.Currency,
			Lock:     asset.Lock,
			Free:     asset.Free,
		})
	}

	return nil
}

func (b *Binance) GetAssets() (types.Assets, error) {
	response := AssetsResponse{}
	err := b.api.Request(api.Request{
		Method: "GET", Endpoint: "/account", Authenticate: true,
	}, nil, &response)

	return types.Assets(response), err
}
