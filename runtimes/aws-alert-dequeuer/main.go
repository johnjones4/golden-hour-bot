package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/johnjones4/golden-hour-bot/lib/runtime"
)

func main() {
	lambda.Start(runtime.AWSAlertDequeuerHandler)
}
