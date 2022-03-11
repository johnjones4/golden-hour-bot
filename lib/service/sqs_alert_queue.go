package service

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/johnjones4/golden-hour-bot/lib/shared"
	"github.com/johnjones4/golden-hour-bot/lib/telegram"
)

type sQSAlertQueueItem struct {
	PredType string          `json:"pType"`
	Reminder shared.Reminder `json:"reminder"`
}

type SQSAlertEnqueuer struct {
	SQS      *sqs.SQS
	QueueURL string
}

func (s *SQSAlertEnqueuer) EnqueueAlerts(predType string, reminders []shared.Reminder) error {
	for _, reminder := range reminders {
		item := sQSAlertQueueItem{predType, reminder}

		bytes, err := json.Marshal(item)
		if err != nil {
			return err
		}

		_, err = s.SQS.SendMessage(&sqs.SendMessageInput{
			QueueUrl:    &s.QueueURL,
			MessageBody: aws.String(string(bytes)),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func ProcessAlert(client telegram.Telegram, messageId, messageBody string) error {
	log.Printf("Proccessing alert message %s", messageId)

	var item sQSAlertQueueItem
	err := json.Unmarshal([]byte(messageBody), &item)
	if err != nil {
		return err
	}

	err = SendAlert(client, item.PredType, item.Reminder)
	if err != nil {
		return err
	}

	return nil
}
