package hitbtc

import (
	"encoding/json"

	"go.sadegh.io/expi/internal/api"
	"go.sadegh.io/expi/models"
)

type AssetsResponse models.Assets

func (r *AssetsResponse) UnmarshalJSON(data []byte) error {
	var v []struct {
		Lock     float64 `json:"reserved,string"`
		Free     float64 `json:"available,string"`
		Currency string  `json:"currency,required"`
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	for _, asset := range v {
		*r = append(*r, models.Asset{
			Currency: asset.Currency,
			Lock:     asset.Lock,
			Free:     asset.Free,
		})
	}

	return nil
}

func (h *HitBTC) GetAssets() (models.Assets, error) {
	response := AssetsResponse{}
	err := h.api.Request(api.Request{
		Method: "GET", Endpoint: "/spot/balance", Authenticate: true,
	}, nil, &response)

	return models.Assets(response), err
}
