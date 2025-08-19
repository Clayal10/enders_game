package cross

import "errors"

var (
	ErrFrameTooSmall      = errors.New("frame too small")
	ErrInvalidMessageType = errors.New("invalid message type")
)
