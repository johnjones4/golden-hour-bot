package shared

import (
	"errors"
	"fmt"
)

type ResponseError struct {
	Parent  error
	Message string
}

func (e ResponseError) Error() string {
	return e.Parent.Error()
}

func ErrorDuplicateReminder(key string, chatId int) ResponseError {
	return ResponseError{
		Parent:  fmt.Errorf("duplicate reminder %s: %d", key, chatId),
		Message: MessageDuplicateReminder,
	}
}

func ErrorTimeNotParsable(str string) ResponseError {
	return ResponseError{
		Parent:  fmt.Errorf("cannot parse time %s", str),
		Message: MessageBadDate,
	}
}

func ErrorNoLocation() ResponseError {
	return ResponseError{
		Parent:  errors.New("no location provided"),
		Message: MessageNoLocation,
	}
}
