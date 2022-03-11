package runtime

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/codingsince1985/geo-golang/openstreetmap"
	"github.com/johnjones4/golden-hour-bot/lib/engine"
	"github.com/johnjones4/golden-hour-bot/lib/service"
	"github.com/johnjones4/golden-hour-bot/lib/telegram"
	"github.com/olebedev/when"
	"github.com/olebedev/when/rules/common"
	"github.com/olebedev/when/rules/en"
)

func AWSWebhookHandler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var update telegram.Update
	err := json.Unmarshal([]byte(event.Body), &update)
	if err != nil {
		return events.APIGatewayProxyResponse{}, logError(err)
	}

	sess, err := session.NewSession()
	if err != nil {
		panic(err)
	}

	db := dynamodb.New(sess)

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

	err = e.ProcessMessage(update.Message)
	if err != nil {
		return events.APIGatewayProxyResponse{}, logError(err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
	}, nil
}
