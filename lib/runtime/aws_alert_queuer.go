package runtime

import (
	"context"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/johnjones4/golden-hour-bot/lib/service"
)

func AWSAlertQueuerHandler(context.Context, events.CloudWatchEvent) error {
	sess, err := session.NewSession()
	if err != nil {
		return logError(err)
	}

	aq := &service.SQSAlertEnqueuer{
		SQS:      sqs.New(sess),
		QueueURL: os.Getenv("SQS_QUEUE_URL"),
	}
	rs := &service.DynamoReminderStorage{
		Db: dynamodb.New(sess),
	}

	err = service.RunAlertCycle(rs, aq)
	if err != nil {
		return logError(err)
	}

	return nil
}
