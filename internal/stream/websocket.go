package stream

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/gorilla/websocket"

	"go.sadegh.io/expi/internal/api"
	"go.sadegh.io/expi/internal/errors"
	"go.sadegh.io/expi/models/event"
)

type Request struct {
	ID           int    `json:"id"`
	Params       any    `json:"params"`
	Method       string `json:"method"`
	Channel      string `json:"ch,omitempty"`
	Endpoint     string `json:"-"`
	Authenticate bool   `json:"-"`
}

type conn struct {
	*websocket.Conn
	exp time.Time
}

func (s *Stream) ListenKey() string {
	if s.listenKey != "" {
		return s.listenKey
	}

	var response map[string]string

	binance := api.NewBinanceAPI(s.publicKey, s.secretKey)
	err := binance.Request(api.Request{
		Method: "POST", Endpoint: "/userDataStream",
	}, nil, &response)
	if err != nil {
		s.events <- &event.Event{
			Topic:    "error",
			Response: fmt.Errorf("authenticate failed: %v", err).Error(),
		}
	}

	s.listenKey = response["listenKey"]

	// Send a PUT request every 30 minutes for keep-alive listen key.
	go func() {
		ticker := time.NewTicker(time.Minute * 30)

		for range ticker.C {
			err = binance.Request(api.Request{
				Method: "PUT", Endpoint: "/userDataStream",
			}, map[string]string{"listenKey": s.listenKey}, nil)
			if err != nil {
				s.events <- &event.Event{
					Topic:    "error",
					Response: fmt.Errorf("reset listen key failed: %v", err).Error(),
				}
			}
		}
	}()

	return s.listenKey
}

func (s *Stream) RandomKey() string {
	var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("_%d", seededRand.Int())
}

func (s *Stream) Request(request Request) error {
	s.requests = append(s.requests, request)
	conn, err := s.getStream(request.Endpoint)
	if err != nil {
		return err
	}

	if request.Authenticate {
		if s.publicKey == "" || s.secretKey == "" {
			return errors.ErrKeysNotSet
		}

		s.authenticate(conn)
	}

	_ = s.limiter.Wait(context.Background())

	if request.Method != "" {
		err = conn.WriteJSON(request)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Stream) Receive(request Request) {
	conn, err := s.getStream(request.Endpoint)
	if err != nil {
		s.events <- &event.Event{
			Topic:    "error",
			Response: err.Error(),
		}
	}

	if s.keepAlive != nil {
		ticker := time.NewTicker(*s.keepAlive)
		go func() {
			for range ticker.C {
				deadline := time.Now().Add(10 * time.Second)
				err := conn.WriteControl(websocket.PongMessage, []byte{}, deadline)
				if err != nil {
					s.events <- &event.Event{
						Topic:    "error",
						Response: fmt.Errorf("write pong message: %v", err).Error(),
					}
				}

				if time.Now().After(conn.exp) {
					if _, err = s.addStream(request.Endpoint); err != nil {
						s.events <- &event.Event{
							Topic:    "error",
							Response: fmt.Errorf("reset stream conn failed: %v", err).Error(),
						}
					}
				}
			}
		}()
	}

	for {
		var (
			msg     []byte
			evt     = &event.Event{}
			payload = map[string]json.RawMessage{}
		)

		_, msg, err = conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				s.events <- &event.Event{
					Topic:    "error",
					Response: fmt.Errorf("read message: %v", err).Error(),
				}
			}

			break
		}

		err = json.Unmarshal(msg, &payload)
		if err != nil {
			s.events <- &event.Event{
				Topic:    "error",
				Response: fmt.Errorf("unmarshal payload: %v", err).Error(),
			}

			continue
		}

		if _, ok := payload["error"]; ok {
			s.events <- &event.Event{
				Response: string(msg),
				Topic:    "error",
			}

			continue
		}

		// Serialize json-rpc2 response message
		if method, ok := payload["method"]; ok {
			evt.Topic = string(method)[1 : len(string(method))-1]
			evt.Response, _ = json.Marshal(payload["params"])
		}

		// Serialize hitbtc subscription
		if method, ok := payload["ch"]; ok {
			evt.Topic = string(method)[1 : len(string(method))-1]
			evt.Response, _ = json.Marshal(payload["update"])
		}

		// Serialize binance subscription
		if method, ok := payload["e"]; ok {
			evt.Topic = string(method)[1 : len(string(method))-1]
			evt.Response, _ = json.Marshal(&payload)
		}

		if evt.Topic != "" {
			s.events <- evt
		}
	}
}

func (s *Stream) receiveEvents() {
	for evt := range s.events {
		s.receive(evt)
	}
}

func (s *Stream) addStream(stream string) (*conn, error) {
	url := s.baseURL
	if strings.Contains(stream, "/") {
		url += stream
	}

	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, err
	}

	conn := &conn{exp: time.Now().Add(time.Hour * 24), Conn: c}
	s.conns[stream] = conn

	for _, request := range s.requests {
		if request.Endpoint != stream {
			continue
		}

		err = s.Request(request)
		if err != nil {
			s.events <- &event.Event{
				Topic:    "error",
				Response: fmt.Errorf("resubscribe stream failed: %v", err).Error(),
			}
		}
	}

	return conn, err
}

func (s *Stream) getStream(stream string) (*conn, error) {
	if conn, ok := s.conns[stream]; ok {
		return conn, nil
	}

	conn, err := s.addStream(stream)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
