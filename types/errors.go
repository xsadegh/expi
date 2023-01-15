package types

import (
	"encoding/json"
	"fmt"
)

var ErrKeysNotSet = fmt.Errorf("api keys not set")

type ApiErr struct {
	Code    int
	Message string
}

func (e *ApiErr) UnmarshalJSON(data []byte) error {
	var v struct {
		Code    int    `json:"code"`
		Message string `json:"msg"`
		Error   *struct {
			Code        int    `json:"code"`
			Message     string `json:"message"`
			Description string `json:"description"`
		} `json:"error"`
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	if v.Error != nil {
		e.Code = v.Error.Code
		if v.Error.Description != "" {
			e.Message = fmt.Sprintf("%s: [%s]", v.Error.Message, v.Error.Description)
		} else {
			e.Message = v.Error.Message
		}
	} else {
		e.Code = v.Code
		e.Message = v.Message
	}

	return nil
}

func (e *ApiErr) Error() string {
	return fmt.Sprintf("[%d]: %s", e.Code, e.Message)
}
