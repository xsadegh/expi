package hitbtc

import (
	"go.sadegh.io/expi/internal/api"
	"go.sadegh.io/expi/internal/stream"
	"go.sadegh.io/expi/models/event"
)

type HitBTC struct {
	api      *api.API
	stream   *stream.Stream
	receiver event.Receiver

	publicKey string
	secretKey string
}

type Option func(*HitBTC)

func NewHitBTC(opts ...Option) *HitBTC {
	hitbtc := &HitBTC{}
	for _, opt := range opts {
		opt(hitbtc)
	}
	hitbtc.api = api.NewHitbtcAPI(hitbtc.publicKey, hitbtc.secretKey)
	if hitbtc.receiver != nil {
		hitbtc.stream = stream.NewHitbtcStream(
			hitbtc.receive, hitbtc.publicKey, hitbtc.secretKey,
		)
	}

	return hitbtc
}

func WithKeys(public, secret string) Option {
	return func(h *HitBTC) {
		h.publicKey, h.secretKey = public, secret
	}
}

func WithReceiver(receiver event.Receiver) Option {
	return func(b *HitBTC) {
		b.receiver = receiver
	}
}
