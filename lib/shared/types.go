package shared

import (
	"time"

	"github.com/johnjones4/golden-hour-bot/lib/telegram"
)

type StateEngine interface {
	GetChatState(id int) (string, interface{}, error)
	SetChatState(id int, state string, info interface{}) error
}

type PredictionRequest struct {
	ChatId         int         `json:"chatId"`
	PredictionType string      `json:"predictionType"`
	LocationString string      `json:"locationString"`
	Location       Coordinates `json:"location"`
	When           string      `json:"when"`
}

type Queue interface {
	EnqueueBasicMessage(telegram.OutgoingMessage) error
	EnqueueReplyKeyboardMarkupMessage(telegram.OutgoingReplyKeyboardMarkupMessage) error
	EnqueueReplyKeyboardRemoveMessage(telegram.OutgoingReplyKeyboardRemoveMessage) error
	EnqueuePrediction(req PredictionRequest) error
	EnqueueReminder(req RemindRequest) error
}

type Coordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type RemindRequest struct {
	ChatId         int         `json:"chatId"`
	LocationString string      `json:"locationString"`
	Location       Coordinates `json:"location"`
}

type LocationDetails struct {
	City  string `json:"city"`
	State string `json:"state"`
}

type Reminder struct {
	ChatId          int             `json:"chatId"`
	Location        Coordinates     `json:"location"`
	LocationDetails LocationDetails `json:"locationDetails"`
	Timezone        string          `json:"timezone"`
}

type Region struct {
	Region           string      `json:"region"`
	Location         Coordinates `json:"location"`
	NextSunriseAlert time.Time   `json:"nextSunriseAlert"`
	NextSunsetAlert  time.Time   `json:"nextSunsetAlert"`
	Timezone         string      `json:"timezone"`
}

type ReminderStorage interface {
	SaveReminder(Reminder) error
	GetRegions() ([]Region, error)
	UpdateRegion(r Region) error
	GetRemindersInRegion(region string) ([]Reminder, error)
}

type AlertQueue interface {
	EnqueueAlerts(predType string, reminders []Reminder) error
}

type ParsedPredictionRequest struct {
	PredictionType  string
	Location        Coordinates
	LocationDetails LocationDetails
	Timezone        string
	Date            time.Time
}
