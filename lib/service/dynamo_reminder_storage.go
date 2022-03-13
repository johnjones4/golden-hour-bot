package service

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/johnjones4/golden-hour-bot/lib/shared"
)

const (
	reminderStorageTable      = "golden-reminders"
	reminderIndexStorageTable = "golden-reminders-index"
)

type ReminderStorageItem struct {
	shared.Reminder
	Region string `json:"region"`
}

type DynamoReminderStorage struct {
	Db *dynamodb.DynamoDB
}

func (rs *DynamoReminderStorage) SaveReminder(r shared.Reminder) error {
	reminderStorage, exists, err := rs.GetReminder(r.GetRegionKey(), r.ChatId)
	if err != nil {
		return err
	}
	if exists {
		return shared.ErrorDuplicateReminder(r.GetRegionKey(), r.ChatId)
	}

	reminderStorage = ReminderStorageItem{
		Region:   r.GetRegionKey(),
		Reminder: r,
	}

	encoded, err := dynamodbattribute.MarshalMap(reminderStorage)
	if err != nil {
		return err
	}

	_, err = rs.Db.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(reminderStorageTable),
		Item:      encoded,
	})
	if err != nil {
		return err
	}

	c := r.GetRegion()
	err = rs.upsertIndex(r.GetRegionKey(), c, r.Timezone)
	if err != nil {
		return err
	}

	return nil
}

func (rs *DynamoReminderStorage) GetReminder(region string, chatId int) (ReminderStorageItem, bool, error) {
	result, err := rs.Db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(reminderStorageTable),
		Key: map[string]*dynamodb.AttributeValue{
			"region": {
				S: aws.String(region),
			},
			"chatId": {
				N: aws.String(fmt.Sprint(chatId)),
			},
		},
	})

	if err != nil {
		if _, ok := err.(*dynamodb.ResourceNotFoundException); ok {
			return ReminderStorageItem{}, false, nil
		}
		return ReminderStorageItem{}, false, err
	}

	if result.Item == nil {
		return ReminderStorageItem{}, false, nil
	}

	var reminderStorage ReminderStorageItem
	err = dynamodbattribute.UnmarshalMap(result.Item, &reminderStorage)
	if err != nil {
		return ReminderStorageItem{}, false, err
	}

	return reminderStorage, false, nil
}

func (rs *DynamoReminderStorage) GetRegions() ([]shared.Region, error) {
	var lastPage map[string]*dynamodb.AttributeValue = nil
	firstPage := true
	aggregator := make([]shared.Region, 0)

	for firstPage || lastPage != nil {
		firstPage = false
		result, err := rs.Db.Scan(&dynamodb.ScanInput{
			TableName:         aws.String(reminderIndexStorageTable),
			ExclusiveStartKey: lastPage,
		})
		if err != nil {
			return nil, err
		}

		var results []shared.Region
		err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &results)
		if err != nil {
			return nil, err
		}

		aggregator = append(aggregator, results...)

		if len(result.LastEvaluatedKey) == 0 {
			break
		}

		lastPage = result.LastEvaluatedKey

	}

	return aggregator, nil
}

func (rs *DynamoReminderStorage) GetRemindersInRegion(region string) ([]shared.Reminder, error) {
	var lastPage map[string]*dynamodb.AttributeValue = nil
	firstPage := true
	aggregator := make([]shared.Reminder, 0)

	for firstPage || lastPage != nil {
		result, err := rs.Db.Query(&dynamodb.QueryInput{
			TableName:              aws.String(reminderStorageTable),
			KeyConditionExpression: aws.String("#r = :r"),
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":r": {
					S: aws.String(region),
				},
			},
			ExpressionAttributeNames: map[string]*string{
				"#r": aws.String("region"),
			},
		})

		if err != nil {
			return nil, err
		}

		var results []shared.Reminder
		err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &results)
		if err != nil {
			return nil, err
		}

		aggregator = append(aggregator, results...)

		if len(result.LastEvaluatedKey) == 0 {
			break
		}

		lastPage = result.LastEvaluatedKey
	}

	return aggregator, nil
}

func (rs *DynamoReminderStorage) UpdateRegion(r shared.Region) error {
	encoded, err := dynamodbattribute.MarshalMap(r)
	if err != nil {
		return err
	}

	_, err = rs.Db.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(reminderIndexStorageTable),
		Item:      encoded,
	})
	if err != nil {
		return err
	}

	return nil
}

func (rs *DynamoReminderStorage) upsertIndex(region string, coord shared.Coordinates, tz string) error {
	result, err := rs.Db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(reminderIndexStorageTable),
		Key: map[string]*dynamodb.AttributeValue{
			"region": {
				S: aws.String(region),
			},
		},
	})

	if err != nil {
		if _, ok := err.(*dynamodb.ResourceNotFoundException); !ok {
			return err
		}
	}

	if result.Item != nil {
		return nil
	}

	item := shared.Region{
		Region:   region,
		Location: coord,
		Timezone: tz,
	}

	encoded, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		return err
	}

	_, err = rs.Db.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(reminderIndexStorageTable),
		Item:      encoded,
	})
	if err != nil {
		return err
	}

	return nil
}
