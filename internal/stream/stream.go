package stream

import (
	"fmt"
	"time"

	"golang.org/x/time/rate"

	"go.sadegh.io/expi/models/event"
)

type Stream struct {
	receive      event.Receiver
	baseURL      string
	publicKey    string
	secretKey    string
	listenKey    string
	authenticate authenticator
	keepAlive    *time.Duration
	requests     []Request
	limiter      *rate.Limiter
	events       chan *event.Event
	conns        map[string]*conn
}

type Option func(*Stream)
type authenticator func(conn *conn)

func NewStream(opts ...Option) *Stream {
	stream := &Stream{
		events: make(chan *event.Event),
		conns:  make(map[string]*conn),
	}
	for _, opt := range opts {
		opt(stream)
	}

	return stream
}

func NewHitbtcStream(receiver event.Receiver, keys ...string) *Stream {
	var options []Option
	if len(keys) == 2 {
		options = append(options, streamWithKeys(keys[0], keys[1]))
	}

	options = append(options, func(s *Stream) {
		s.authenticate = func(conn *conn) {
			params := struct {
				SecretKey string `json:"secret_key,required"`
				PublicKey string `json:"api_key,required"`
				Algorithm string `json:"type,required"`
			}{
				SecretKey: s.secretKey,
				PublicKey: s.publicKey,
				Algorithm: "BASIC",
			}

			err := conn.WriteJSON(Request{Method: "login", Params: params})
			if err != nil {
				s.events <- &event.Event{
					Topic:    "error",
					Response: fmt.Errorf("authenticate failed: %v", err).Error(),
				}
			}
		}
	})

	options = append(options, streamWithReceiver(receiver), streamWithRateLimit(20))

	stream := NewStream(options...)
	stream.baseURL = "wss://api.hitbtc.com/api/3/ws"
	go stream.receiveEvents()

	return stream
}

func NewBinanceStream(receiver event.Receiver, keys ...string) *Stream {
	var options []Option
	if len(keys) == 2 {
		options = append(options, streamWithKeys(keys[0], keys[1]))
	}

	options = append(options,
		streamWithReceiver(receiver),
		streamWithRateLimit(5),
		streamWithKeepAlive(3*time.Minute),
	)

	stream := NewStream(options...)
	stream.baseURL = "wss://stream.binance.com:9443/ws"
	go stream.receiveEvents()

	return stream
}

func streamWithKeepAlive(keepAlive time.Duration) Option {
	return func(s *Stream) {
		s.keepAlive = &keepAlive
	}
}

func streamWithRateLimit(burst int) Option {
	return func(s *Stream) {
		s.limiter = rate.NewLimiter(rate.Every(time.Second), burst)
	}
}

func streamWithReceiver(receiver event.Receiver) Option {
	return func(s *Stream) {
		s.receive = receiver
	}
}

func streamWithKeys(public, secret string) Option {
	return func(s *Stream) {
		s.publicKey, s.secretKey = public, secret
	}
}
