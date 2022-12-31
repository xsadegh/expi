package binance

import (
	"encoding/json"

	"go.sadegh.io/expi/internal/api"
	"go.sadegh.io/expi/models"
)

type AssetsResponse models.Assets

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
		*r = append(*r, models.Asset{
			Currency: asset.Currency,
			Lock:     asset.Lock,
			Free:     asset.Free,
		})
	}

	return nil
}

func (b *Binance) GetAssets() (models.Assets, error) {
	response := AssetsResponse{}
	err := b.api.Request(api.Request{
		Method: "GET", Endpoint: "/account", Authenticate: true,
	}, nil, &response)

	return models.Assets(response), err
}
