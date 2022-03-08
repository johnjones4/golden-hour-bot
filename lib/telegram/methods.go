package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Telegram struct {
	Token string
}

func (t Telegram) callMethod(name string, parameters interface{}, response interface{}) error {
	bodyBytes := []byte{}
	var err error

	if parameters != nil {
		bodyBytes, err = json.Marshal(parameters)
		if err != nil {
			return err
		}
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("https://api.telegram.org/%s/%s", t.Token, name), bytes.NewBuffer(bodyBytes))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(responseBody, response)
	if err != nil {
		return err
	}

	return nil
}

func (t Telegram) SendMessage(m OutgoingMessage) (SendMessageResponse, error) {
	var resp SendMessageResponse
	err := t.callMethod("sendMessage", nil, &resp)
	if err != nil {
		return SendMessageResponse{}, err
	}
	return resp, nil
}
