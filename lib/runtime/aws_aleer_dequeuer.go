package runtime

import (
	"context"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/johnjones4/golden-hour-bot/lib/service"
	"github.com/johnjones4/golden-hour-bot/lib/telegram"
)

func AWSAlertDequeuerHandler(ctx context.Context, event events.SQSEvent) error {
	tClient := telegram.Telegram{
		Token: os.Getenv("TELEGRAM_TOKEN"),
	}

	for _, message := range event.Records {
		err := service.ProcessAlert(tClient, message.MessageId, message.Body)
		if err != nil {
			return logError(err)
		}
	}

	return nil
}
