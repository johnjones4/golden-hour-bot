package engine

import (
	"log"
	"strings"

	"github.com/johnjones4/golden-hour-bot/lib/shared"
	"github.com/johnjones4/golden-hour-bot/lib/telegram"
)

var (
	commandStart      = "/start"
	commandGetSunrise = "/" + shared.PredictionTypeSunrise
	commandGetSunset  = "/" + shared.PredictionTypeSunset
	commandRemind     = "/remind"
)

type Engine struct {
	StateEngine shared.StateEngine
	Queue       shared.Queue
}

func (e *Engine) ProcessMessage(message telegram.IncomingMessage) error {
	log.Println(message)

	state, uncastInfo, err := e.StateEngine.GetChatState(message.Chat.Id)
	if err != nil {
		return err
	}

	cmd, foundCmd := message.GetCommand()

	switch state {

	case shared.StateIdle:
		if foundCmd {
			switch cmd.Command {
			case commandStart:
				return e.Queue.EnqueueBasicMessage(telegram.OutgoingMessage{
					ChatId: message.Chat.Id,
					Message: telegram.Message{
						Text: shared.MessageWelcome,
					},
				})
			case commandGetSunrise, commandGetSunset:
				req := shared.PredictionRequest{
					PredictionType: cmd.Command[1:],
					ChatId:         message.Chat.Id,
				}

				if strings.TrimSpace(cmd.Extra) != "" {
					req.When = strings.TrimSpace(cmd.Extra)
				}

				if !message.Location.IsZero() {
					req.Location = shared.Coordinates{
						Latitude:  message.Location.Latitude,
						Longitude: message.Location.Longitude,
					}
				}

				return checkAndHandleImmediate(e.Queue, e.StateEngine, req)
			case commandRemind:
				req := shared.RemindRequest{
					ChatId: message.Chat.Id,
				}

				if strings.TrimSpace(cmd.Extra) != "" {
					req.LocationString = strings.TrimSpace(cmd.Extra)
				}

				return checkAndHandleRemind(e.Queue, e.StateEngine, req)
			}
		}
		return e.Queue.EnqueueBasicMessage(telegram.OutgoingMessage{
			ChatId: message.Chat.Id,
			Message: telegram.Message{
				Text: shared.MessageConfused,
			},
		})

	case shared.StateImmediateNeedLocation:
		if req, ok := uncastInfo.(shared.PredictionRequest); ok {
			if !message.Location.IsZero() {
				req.Location = shared.Coordinates{
					Latitude:  message.Location.Latitude,
					Longitude: message.Location.Longitude,
				}
			} else if strings.TrimSpace(message.Text) != "" {
				req.LocationString = strings.TrimSpace(message.Text)
			} else {
				return goToImmediateNeedLocation(e.Queue, e.StateEngine, req)
			}

			err = clearLocationAsk(e.Queue, message.Chat.Id)
			if err != nil {
				return err
			}

			return checkAndHandleImmediate(e.Queue, e.StateEngine, req)
		}

	case shared.StateImmediateNeedDate:
		if req, ok := uncastInfo.(shared.PredictionRequest); ok {
			if strings.TrimSpace(message.Text) != "" {
				req.When = strings.TrimSpace(message.Text)
			}

			return checkAndHandleImmediate(e.Queue, e.StateEngine, req)
		}

	case shared.StateRemindNeedLocation:
		if req, ok := uncastInfo.(shared.RemindRequest); ok {
			if !message.Location.IsZero() {
				req.Location = shared.Coordinates{
					Latitude:  message.Location.Latitude,
					Longitude: message.Location.Longitude,
				}
			} else if strings.TrimSpace(message.Text) != "" {
				req.LocationString = strings.TrimSpace(message.Text)
			} else {
				return goToRemindNeedLocation(e.Queue, e.StateEngine, req)
			}

			err = clearLocationAsk(e.Queue, message.Chat.Id)
			if err != nil {
				return err
			}

			return checkAndHandleRemind(e.Queue, e.StateEngine, req)
		}

	}

	return nil
}
