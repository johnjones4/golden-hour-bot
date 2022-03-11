package runtime

import (
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/codingsince1985/geo-golang/openstreetmap"
	"github.com/johnjones4/golden-hour-bot/lib/engine"
	"github.com/johnjones4/golden-hour-bot/lib/service"
	"github.com/johnjones4/golden-hour-bot/lib/telegram"
	"github.com/olebedev/when"
	"github.com/olebedev/when/rules/common"
	"github.com/olebedev/when/rules/en"
)

func StartLocalServerRuntime() {
	sess, err := session.NewSession()
	if err != nil {
		panic(err)
	}

	db := makeLocalDynamoConnection(sess)

	w := when.New(nil)
	w.Add(en.All...)
	w.Add(common.All...)

	geocoder := openstreetmap.Geocoder()

	mq := &service.DirectQueue{
		Client: telegram.Telegram{
			Token: os.Getenv("TELEGRAM_TOKEN"),
		},
		GeoNames: service.GeoNames{
			Username: os.Getenv("GEONAMES_USERNAME"),
		},
		PredictionParser: service.PredictionRequestParser{
			GeoNames: service.GeoNames{
				Username: os.Getenv("GEONAMES_USERNAME"),
			},
			DateParser: w,
			Geocoder:   geocoder,
		},
		ReminderStorage: &service.DynamoReminderStorage{
			Db: db,
		},
		Geocoder: geocoder,
	}

	e := engine.Engine{
		StateEngine: &service.DynamoStateEngine{
			Db: db,
		},
		Queue: mq,
	}

	runLocalServer(&e)
}

func StartLocalAlertQueuer() {
	sess, err := session.NewSession()
	if err != nil {
		panic(err)
	}

	aq := &service.SQSAlertEnqueuer{
		SQS:      makeLocalSQSConection(sess),
		QueueURL: os.Getenv("SQS_QUEUE_URL"),
	}
	rs := &service.DynamoReminderStorage{
		Db: makeLocalDynamoConnection(sess),
	}
	for {
		err := service.RunAlertCycle(rs, aq)
		if err != nil {
			panic(err)
		}

		time.Sleep(time.Second * 30)
	}
}

func StartLocalDequeuerRuntime() {
	sess, err := session.NewSession()
	if err != nil {
		panic(err)
	}

	tClient := telegram.Telegram{
		Token: os.Getenv("TELEGRAM_TOKEN"),
	}

	sqsClient := makeLocalSQSConection(sess)

	for {
		messages, err := sqsClient.ReceiveMessage(&sqs.ReceiveMessageInput{
			QueueUrl: aws.String(os.Getenv("SQS_QUEUE_URL")),
		})
		if err != nil {
			panic(err)
		}
		for _, message := range messages.Messages {
			err = service.ProcessAlert(tClient, *message.MessageId, *message.Body)
			if err != nil {
				panic(err)
			}
		}
		time.Sleep(time.Second)
	}
}

func makeLocalDynamoConnection(sess *session.Session) *dynamodb.DynamoDB {
	endpoint := os.Getenv("DYNAMO_ENDPOINT")
	cfg := aws.Config{}
	cfg.Endpoint = aws.String(endpoint)
	cfg.Region = aws.String("us-east-1")
	cfg.Credentials = credentials.NewStaticCredentials("dev", "dev", "")
	return dynamodb.New(sess, &cfg)
}

func makeLocalSQSConection(sess *session.Session) *sqs.SQS {
	endpoint := os.Getenv("SQS_ENDPOINT")
	cfg := aws.Config{}
	cfg.Endpoint = aws.String(endpoint)
	cfg.Region = aws.String("us-east-1")
	cfg.Credentials = credentials.NewStaticCredentials("dev", "dev", "")
	return sqs.New(sess, &cfg)
}
