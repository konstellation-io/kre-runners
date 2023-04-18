package errors

import (
	"errors"
	"fmt"
)

var ErrUndefinedObjectStore = errors.New("the object store does not exist")
var ErrMessageToBig = errors.New("compressed message exceeds maximum size allowed")
var ErrMsgAck = "Error in message ack: %s"
var ErrEmptyPayload = errors.New("the payload cannot be empty")

// Wrapper creates a function that returns errors starts with a given message.
func Wrapper(message string) func(params ...interface{}) error {
	return func(params ...interface{}) error {
		return fmt.Errorf(message, params...)
	}
}
