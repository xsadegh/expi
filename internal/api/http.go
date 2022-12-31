package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"go.sadegh.io/expi/internal/errors"
)

type Request struct {
	Method       string
	Endpoint     string
	Authenticate bool
}

func (h *API) Request(request Request, params, response any) (err error) {
	parsedURL, _ := url.ParseRequestURI(h.baseURL)
	parsedURL.Path = parsedURL.Path + request.Endpoint

	bin, err := h.parseParams(parsedURL, params)
	if err != nil {
		return
	}

	if request.Authenticate {
		if h.publicKey == "" || h.secretKey == "" {
			return errors.ErrKeysNotSet
		}
	}

	_ = h.limiter.Wait(context.Background())

	if request.Authenticate && h.authOnQuery != nil {
		h.authOnQuery(parsedURL, h.publicKey, h.secretKey)
	} else {
		parsedURL.RawQuery = parsedURL.Query().Encode()
	}

	var req *http.Request
	if request.Authenticate {
		req, err = http.NewRequest(request.Method, parsedURL.String(), nil)
	} else {
		req, err = http.NewRequest(request.Method, parsedURL.String(), bytes.NewBuffer(bin))
	}

	if err != nil {
		return err
	}

	if request.Authenticate && h.authOnRequest != nil {
		h.authOnRequest(req, h.publicKey, h.secretKey)
	}

	for header, value := range h.headers {
		req.Header.Add(header, value)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != 200 {
		apiErr := &errors.ApiErr{}
		if err = json.NewDecoder(res.Body).Decode(&apiErr); err != nil {
			return
		}

		return apiErr
	}

	if err = json.NewDecoder(res.Body).Decode(&response); err != nil {
		return
	}

	return
}

func (h *API) parseParams(base *url.URL, params any) (bin []byte, err error) {
	if params == nil {
		return
	}

	bin, err = json.Marshal(params)
	if err != nil {
		return nil, err
	}
	m := map[string]interface{}{}

	err = json.Unmarshal(bin, &m)
	if err != nil {
		return nil, err
	}

	q := base.Query()
	for k, v := range m {
		q.Add(k, fmt.Sprintf("%v", v))
	}

	base.RawQuery = q.Encode()

	return bin, err
}
