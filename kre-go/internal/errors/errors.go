package errors

import "errors"

var ErrUndefinedObjectStore = errors.New("the object store does not exist")
var ErrMessageToBig = errors.New("compressed message exceeds maximum size allowed of 1 MB")
var ErrMsgAck = "Error in message ack: %s"
var ErrEmptyPayload = errors.New("the payload cannot be empty")
