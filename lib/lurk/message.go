package lurk

import (
	"encoding/binary"

	"github.com/Clayal10/enders_game/lib/cross"
)

type messageType byte

// Exported message types.
const (
	TypeMessage    messageType = 1
	TypeChangeRoom messageType = 2
	TypeFight      messageType = 3
	TypePVPFight   messageType = 4
	TypeLoot       messageType = 5
	TypeStart      messageType = 6
	TypeError      messageType = 7
	TypeAccept     messageType = 8
	TypeRoom       messageType = 9
	TypeCharacter  messageType = 10
	TypeGame       messageType = 11
	TypeLeave      messageType = 12
	TypeConnection messageType = 13
	TypeVersion    messageType = 14
)

// Those using this library will need to use this function on the returned interface
// to know which type will need to be used for assertion.
type LurkMessage interface {
	GetType() messageType
}

// Unmarshal takes a slice of bytes and returns a LurkMessage interface object. The
// LurkMessage then needs to be type asserted based on the type returned from the
// 'GetType()' function.
func Unmarshal(data []byte) (LurkMessage, error) {
	if err := validate(data); err != nil {
		return nil, err
	}

	// Various unmarshaling and returning of their respective types.
	switch messageType(data[0]) {
	case TypeMessage:
		return unmarshalMessage(data)
	case TypeChangeRoom:
	case TypeFight:
	case TypePVPFight:
	case TypeLoot:
	case TypeStart:
	case TypeError:
	case TypeAccept:
	case TypeRoom:
	case TypeCharacter:
	case TypeGame:
	case TypeLeave:
	case TypeConnection:
	case TypeVersion:
	}
	return nil, cross.ErrInvalidMessageType
}

// Marshal Will take any LurkMessage object and return a byte array
// ready for messaging.
func Marshal(lm any) ([]byte, error) {
	msg, ok := lm.(LurkMessage)
	if !ok {
		return nil, cross.ErrInvalidMessageType
	}
	switch msg.GetType() {
	case TypeMessage:
	case TypeChangeRoom:
	case TypeFight:
	case TypePVPFight:
	case TypeLoot:
	case TypeStart:
	case TypeError:
	case TypeAccept:
	case TypeRoom:
	case TypeCharacter:
	case TypeGame:
	case TypeLeave:
	case TypeConnection:
	case TypeVersion:
	}
	return nil, cross.ErrInvalidMessageType
}

func validate(data []byte) error {
	if len(data) < 1 {
		return cross.ErrFrameTooSmall
	}
	if data[0] < 1 || data[0] > 14 {
		return cross.ErrInvalidMessageType
	}
	return nil
}

type Message struct {
	Type      messageType
	RName     string // max 32 bytes. All fields noted with bytes are null terminated '\x00'.
	SName     string // max 30 bytes
	Text      string
	Narration bool
}

func (m *Message) GetType() messageType {
	return m.Type
}

func unmarshalMessage(data []byte) (*Message, error) {
	if len(data) < 8 {
		return nil, cross.ErrFrameTooSmall
	}
	m := &Message{}
	offset := 0
	m.Type = messageType(data[offset])
	offset++
	msgLen := binary.LittleEndian.Uint16(data[offset:])
	offset += 2

	l := getNulTermLen(data[offset:])
	m.RName = string(data[offset : offset+l])
	offset += l + 1
	l = getNulTermLen(data[offset:])
	m.SName = string(data[offset : offset+l])
	offset += l + 1

	// Check for narration
	if data[offset] == 1 {
		m.Narration = true
	}
	offset++
	m.Text = string(data[offset : offset+int(msgLen)])

	return m, nil
}

type ChangeRoom struct {
	Type       messageType
	RoomNumber int
}

func (cr *ChangeRoom) GetType() messageType {
	return cr.Type
}

type Fight struct {
	Type messageType
}

func (f *Fight) GetType() messageType {
	return f.Type
}

type PVPFight struct {
	Type       messageType
	TargetName string // 32 bytes
}

func (pvp *PVPFight) GetType() messageType {
	return pvp.Type
}

type Loot struct {
	Type       messageType
	TargetName string // 32 bytes
}

func (l *Loot) GetType() messageType {
	return l.Type
}

type Start struct {
	Type messageType
}

func (s *Start) GetType() messageType {
	return s.Type
}

type Error struct {
	Type       messageType
	ErrCode    cross.ErrCode
	ErrMessage string
}

func (e *Error) GetType() messageType {
	return e.Type
}

type Accept struct {
	Type   messageType
	Action messageType
}

func (a *Accept) GetType() messageType {
	return a.Type
}

type Room struct {
	Type       messageType
	RoomNumber uint16
	RoomName   string // 32 bytes
	RoomDesc   string
}

func (r *Room) GetType() messageType {
	return r.Type
}

type Character struct {
	Type       messageType
	Name       string // 32 bytes
	Flags      map[string]bool
	Attack     uint16
	Defense    uint16
	Regen      uint16
	Health     int16
	Gold       uint16
	RoomNum    uint16
	PlayerDesc string
}

func (c *Character) GetType() messageType {
	return c.Type
}

type Game struct {
	Type          messageType
	InitialPoints uint16
	StatLimit     uint16
	GameDesc      string
}

func (g *Game) GetType() messageType {
	return g.Type
}

type Leave struct {
	Type messageType
}

func (l *Leave) GetType() messageType {
	return l.Type
}

type Connection struct {
	Type       messageType
	RoomNumber uint16
	RoomName   string //32 bytes
	RoomDesc   string
}

func (c *Connection) GetType() messageType {
	return c.Type
}

type Version struct {
	Type       messageType
	Major      byte
	Minor      byte
	Extensions [][]byte // For now. Turn into object when we know what it is.
}

func (v *Version) GetType() messageType {
	return v.Type
}

const maxStringLen = 32

// data should be a slice starting at the start of a null terminated string.
func getNulTermLen(data []byte) (length int) {
	for _, b := range data {
		if b == '\x00' || length == maxStringLen {
			break
		}
		length++
	}
	return
}
