package binance

import (
	"go.sadegh.io/expi/internal/api"
	"go.sadegh.io/expi/internal/stream"
	"go.sadegh.io/expi/types/event"
)

type Binance struct {
	api      *api.API
	stream   *stream.Stream
	receiver event.Receiver

	publicKey string
	secretKey string
}

type Option func(*Binance)

func NewBinance(opts ...Option) *Binance {
	binance := &Binance{}
	for _, opt := range opts {
		opt(binance)
	}
	binance.api = api.NewBinanceAPI(binance.publicKey, binance.secretKey)
	if binance.receiver != nil {
		binance.stream = stream.NewBinanceStream(
			binance.receive, binance.publicKey, binance.secretKey,
		)
	}

	return binance
}

func WithKeys(public, secret string) Option {
	return func(b *Binance) {
		b.publicKey, b.secretKey = public, secret
	}
}

func WithReceiver(receiver event.Receiver) Option {
	return func(b *Binance) {
		b.receiver = receiver
	}
}
