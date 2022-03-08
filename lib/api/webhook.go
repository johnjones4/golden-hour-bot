package api

import (
	"encoding/json"
	"io"
	"lib/telegram"
	"log"
	"net/http"
)

func handleError(resp http.ResponseWriter, err error) {
	log.Println(err)
	resp.WriteHeader(http.StatusInternalServerError)
}

func MakeWebhook(tClient telegram.Telegram) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		bodybytes, err := io.ReadAll(req.Body)
		if err != nil {
			handleError(resp, err)
			return
		}

		var update telegram.Update
		err = json.Unmarshal(bodybytes, &update)
		if err != nil {
			handleError(resp, err)
			return
		}

		log.Println(update)

		resp.WriteHeader(200)
	})
}
