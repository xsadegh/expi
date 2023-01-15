package api

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/time/rate"
)

type API struct {
	baseURL       string
	limiter       *rate.Limiter
	headers       map[string]string
	authOnQuery   func(*url.URL, string, string)
	authOnRequest func(*http.Request, string, string)

	publicKey, secretKey string
}

type Option func(*API)

func NewAPI(opts ...Option) *API {
	api := &API{}
	for _, opt := range opts {
		opt(api)
	}

	return api
}

func NewHitbtcAPI(keys ...string) *API {
	var options []Option
	headers := map[string]string{"Content-Type": "application/json"}
	if len(keys) == 2 {
		options = append(options, apiWithKeys(keys[0], keys[1]))
	}
	options = append(options, apiWithHeaders(headers), apiWithRateLimit(30))

	api := NewAPI(options...)
	api.authOnRequest = func(req *http.Request, public, secret string) {
		req.SetBasicAuth(public, secret)
	}
	api.baseURL = "https://api.hitbtc.com/api/3"

	return api
}

func NewBinanceAPI(keys ...string) *API {
	var options []Option
	headers := map[string]string{"Content-type": "application/x-www-form-urlencoded"}
	if len(keys) == 2 {
		headers["X-MBX-APIKEY"] = keys[0]
		options = append(options, apiWithKeys(keys[0], keys[1]))
	}
	options = append(options, apiWithHeaders(headers), apiWithRateLimit(10))

	api := NewAPI(options...)
	api.authOnQuery = func(parsedURL *url.URL, public, secret string) {
		q := parsedURL.Query()
		q.Add("recvWindow", "5000")
		// Timestamp is mandatory in signed request
		q.Add("timestamp", fmt.Sprintf("%v", time.Now().Unix()*1000))
		// Signature needs to be at the last param
		parsedURL.RawQuery = q.Encode() + "&signature=" + generateSignature(secret, q)
	}
	api.baseURL = "https://api.binance.com/api/v3"

	return api
}

func apiWithRateLimit(burst int) Option {
	return func(h *API) {
		h.limiter = rate.NewLimiter(rate.Every(time.Second), burst)
	}
}

func apiWithHeaders(headers map[string]string) Option {
	return func(h *API) {
		h.headers = headers
	}
}

func apiWithKeys(public, secret string) Option {
	return func(h *API) {
		h.publicKey, h.secretKey = public, secret
	}
}

func generateSignature(key string, q url.Values) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(q.Encode()))

	expectedMAC := mac.Sum(nil)

	return hex.EncodeToString(expectedMAC)
}
