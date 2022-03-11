package service

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/johnjones4/golden-hour-bot/lib/shared"
)

const (
	stateEngineTable = "golden-state"
)

type dState struct {
	ChatId   int    `json:"chatId"`
	State    string `json:"state"`
	InfoType string `json:"infoType"`
	Info     string `json:"info"`
}

type DynamoStateEngine struct {
	Db *dynamodb.DynamoDB
}

func (e *DynamoStateEngine) GetChatState(chatId int) (string, interface{}, error) {
	result, err := e.Db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(stateEngineTable),
		Key: map[string]*dynamodb.AttributeValue{
			"chatId": {
				N: aws.String(fmt.Sprint(chatId)),
			},
		},
	})
	if err != nil {
		if _, ok := err.(*dynamodb.ResourceNotFoundException); ok {
			return shared.DefaultState, nil, nil
		}
		return "", nil, err
	}

	if result.Item == nil {
		return shared.DefaultState, nil, nil
	}

	var state dState
	err = dynamodbattribute.UnmarshalMap(result.Item, &state)
	if err != nil {
		return "", nil, err
	}

	switch state.InfoType {
	case "PredictionRequest":
		var req shared.PredictionRequest
		err = json.Unmarshal([]byte(state.Info), &req)
		if err != nil {
			return "", nil, err
		}
		return state.State, req, nil
	case "RemindRequest":
		var req shared.RemindRequest
		err = json.Unmarshal([]byte(state.Info), &req)
		if err != nil {
			return "", nil, err
		}
		return state.State, req, nil
	default:
		return state.State, nil, nil
	}
}

func (e *DynamoStateEngine) SetChatState(id int, state string, info interface{}) error {
	infoStr := ""
	infoType := ""

	if info != nil {
		infoBytes, err := json.Marshal(info)
		if err != nil {
			return nil
		}
		infoStr = string(infoBytes)
		infoType = reflect.TypeOf(info).Name()
	}

	stateOb := dState{
		ChatId:   id,
		State:    state,
		InfoType: infoType,
		Info:     infoStr,
	}

	encoded, err := dynamodbattribute.MarshalMap(stateOb)
	if err != nil {
		return err
	}

	_, err = e.Db.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(stateEngineTable),
		Item:      encoded,
	})

	return err
}
