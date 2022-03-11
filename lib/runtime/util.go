package runtime

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/johnjones4/golden-hour-bot/lib/engine"
	"github.com/johnjones4/golden-hour-bot/lib/telegram"
)

func runLocalServer(e *engine.Engine) {
	http.ListenAndServe(":8080", http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
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

		err = e.ProcessMessage(update.Message)
		if err != nil {
			handleError(resp, err)
			return
		}

		resp.WriteHeader(http.StatusOK)
	}))
}

func handleError(resp http.ResponseWriter, err error) {
	log.Println(err)
	resp.WriteHeader(http.StatusOK)
}

func logError(err error) error {
	log.Printf("[APP ERROR]: %s", err.Error())
	return err
}
