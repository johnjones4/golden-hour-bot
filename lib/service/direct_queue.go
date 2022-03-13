package service

import (
	"log"

	"github.com/codingsince1985/geo-golang"
	"github.com/johnjones4/golden-hour-bot/lib/shared"
	"github.com/johnjones4/golden-hour-bot/lib/telegram"
)

type DirectQueue struct {
	Client           telegram.Telegram
	PredictionParser PredictionRequestParser
	ReminderStorage  shared.ReminderStorage
	Geocoder         geo.Geocoder
}

func (mq *DirectQueue) EnqueueBasicMessage(m1 telegram.OutgoingMessage) error {
	log.Printf("Sending message to %d", m1.ChatId)
	_, err := mq.Client.SendMessage(m1)
	return err
}

func (mq *DirectQueue) EnqueueReplyKeyboardMarkupMessage(m2 telegram.OutgoingReplyKeyboardMarkupMessage) error {
	log.Printf("Sending reply keyboard markup message to %d", m2.ChatId)
	_, err := mq.Client.SendReplyKeyboardMarkupMessage(m2)
	return err
}

func (mq *DirectQueue) EnqueueReplyKeyboardRemoveMessage(m3 telegram.OutgoingReplyKeyboardRemoveMessage) error {
	log.Printf("Sending reply keyboard remove message to %d", m3.ChatId)
	_, err := mq.Client.SendReplyKeyboardRemoveMessage(m3)
	return err
}

func (mq *DirectQueue) EnqueuePrediction(req shared.PredictionRequest) error {
	log.Printf("Creating prediction for %d", req.ChatId)
	predictionRequest, err := mq.PredictionParser.NewParsedPredictionRequest(req)
	if err != nil {
		err = wrapErrorWithMessageCatch(mq.Client, req.ChatId, err)
		if err != nil {
			return nil
		}
	}
	log.Println(predictionRequest)

	messageText, err := MakePrediction(predictionRequest)
	if err != nil {
		err = wrapErrorWithMessageCatch(mq.Client, req.ChatId, err)
		if err != nil {
			return nil
		}
	}

	log.Println(messageText)

	_, err = mq.Client.SendMessage(telegram.OutgoingMessage{
		ChatId: req.ChatId,
		Message: telegram.Message{
			Text: messageText,
		},
	})

	return err
}

func (mq *DirectQueue) EnqueueReminder(req shared.RemindRequest) error {
	log.Printf("Creating reminder for %d", req.ChatId)
	reminder, err := NewReminder(mq.Geocoder, req)
	if err != nil {
		return wrapErrorWithMessageCatch(mq.Client, req.ChatId, err)
	}

	err = mq.ReminderStorage.SaveReminder(reminder)
	if err != nil {
		return wrapErrorWithMessageCatch(mq.Client, req.ChatId, err)
	}

	err = mq.EnqueueBasicMessage(telegram.OutgoingMessage{
		ChatId: req.ChatId,
		Message: telegram.Message{
			Text: shared.MessageReminderSet,
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func wrapErrorWithMessageCatch(tClient telegram.Telegram, chatId int, err error) error {
	if err != nil {
		if mErr, ok := err.(shared.ResponseError); ok {
			log.Printf("Experienced messagable error for %d", chatId)
			_, err = tClient.SendMessage(telegram.OutgoingMessage{
				ChatId: chatId,
				Message: telegram.Message{
					Text: mErr.Message,
				},
			})
			return err
		}
		return err
	}
	return nil
}
