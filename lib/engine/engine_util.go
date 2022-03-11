package engine

import (
	"github.com/johnjones4/golden-hour-bot/lib/shared"
	"github.com/johnjones4/golden-hour-bot/lib/telegram"
)

func askForLocation(mq shared.Queue, chatId int) error {
	return mq.EnqueueReplyKeyboardMarkupMessage(telegram.OutgoingReplyKeyboardMarkupMessage{
		OutgoingMessage: telegram.OutgoingMessage{
			ChatId: chatId,
			Message: telegram.Message{
				Text: shared.MessageShareLocation,
			},
		},
		ReplyMarkup: telegram.ReplyKeyboardMarkup{
			Keyboard: [][]telegram.KeyboardButton{
				{{
					RequestLocation: true,
					Text:            shared.ButtonShareLocation,
				}},
			},
		},
	})
}

func clearLocationAsk(mq shared.Queue, chatId int) error {
	return mq.EnqueueReplyKeyboardRemoveMessage(telegram.OutgoingReplyKeyboardRemoveMessage{
		OutgoingMessage: telegram.OutgoingMessage{
			ChatId: chatId,
			Message: telegram.Message{
				Text: shared.MessageLocationThanks,
			},
		},
		ReplyMarkup: telegram.ReplyKeyboardRemove{
			RemoveKeyboard: true,
		},
	})
}

func goToImmediateNeedLocation(mq shared.Queue, se shared.StateEngine, req shared.PredictionRequest) error {
	err := askForLocation(mq, req.ChatId)
	if err != nil {
		return err
	}
	return se.SetChatState(req.ChatId, shared.StateImmediateNeedLocation, req)
}

func goToRemindNeedLocation(mq shared.Queue, se shared.StateEngine, req shared.RemindRequest) error {
	err := askForLocation(mq, req.ChatId)
	if err != nil {
		return err
	}
	return se.SetChatState(req.ChatId, shared.StateRemindNeedLocation, req)
}

func goToImmediateNeedDate(mq shared.Queue, se shared.StateEngine, req shared.PredictionRequest) error {
	err := mq.EnqueueBasicMessage(telegram.OutgoingMessage{
		ChatId: req.ChatId,
		Message: telegram.Message{
			Text: shared.MessageShareTime,
		},
	})
	if err != nil {
		return err
	}
	return se.SetChatState(req.ChatId, shared.StateImmediateNeedDate, req)
}

func checkAndHandleImmediate(mq shared.Queue, se shared.StateEngine, req shared.PredictionRequest) error {
	if req.When == "" {
		return goToImmediateNeedDate(mq, se, req)
	}

	if req.Location.IsZero() && req.LocationString == "" {
		return goToImmediateNeedLocation(mq, se, req)
	}

	err := mq.EnqueuePrediction(req)
	if err != nil {
		return err
	}

	return se.SetChatState(req.ChatId, shared.StateIdle, nil)
}

func checkAndHandleRemind(mq shared.Queue, se shared.StateEngine, req shared.RemindRequest) error {
	if req.Location.IsZero() && req.LocationString == "" {
		return goToRemindNeedLocation(mq, se, req)
	}

	err := mq.EnqueueReminder(req)
	if err != nil {
		return err
	}

	return se.SetChatState(req.ChatId, shared.StateIdle, nil)
}
