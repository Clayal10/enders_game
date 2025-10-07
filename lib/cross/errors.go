package cross

import (
	"errors"
)

var (
	ErrFrameTooSmall      = errors.New("frame too small")
	ErrInvalidMessageType = errors.New("invalid message type")
	ErrInvalidErrCode     = errors.New("invalid error code")
	ErrNoVariableLength   = errors.New("message does not contain a variable length field")
	ErrUserNotInServer    = errors.New("user does not exist in server")
	ErrInvalidRoomNumber  = errors.New("room number does not exist")
	ErrRoomsNotConnected  = errors.New("rooms are not connected")
	ErrQueueEmpty         = errors.New("queue empty")
)

type ErrCode byte

const (
	Other               ErrCode = 0
	BadRoom             ErrCode = 1
	PlayerAlreadyExists ErrCode = 2
	BadMonster          ErrCode = 3
	StatError           ErrCode = 4
	NotReady            ErrCode = 5
	NoTarget            ErrCode = 6
	NoFight             ErrCode = 7
	NoPVP               ErrCode = 8
	NoError             ErrCode = 255
)
