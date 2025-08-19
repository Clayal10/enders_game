package lurk

import "github.com/Clayal10/enders_game/lib/cross"

type messageType byte

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

type LurkMessage interface {
	GetType() messageType
	// Give it a function that can perform the action they need?
}

// Unmarshal will take a raw frame and turn it into the appropriate type that
// satisfies the LurkMessage interface.
func Unmarshal(data []byte) (LurkMessage, error) {
	msg, err := validate(data)
	if err != nil {
		return nil, err
	}

	// Various unmarshaling and returning of their respective types.
	switch msg {
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

// Marshal Will take any LurkMessage object and return a byte array
// ready for messaging.
func Marshal(lm LurkMessage) []byte

func validate(data []byte) (messageType, error) {
	if len(data) < 1 {
		return 0, cross.ErrFrameTooSmall
	}
	if data[0] < 1 || data[0] > 14 {
		return 0, cross.ErrInvalidMessageType
	}
	return messageType(data[0]), nil
}

type Message struct {
	Type messageType
}

func (m *Message) GetType() messageType {
	return m.Type
}

type ChangeRoom struct {
	Type messageType
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
	Type messageType
}

func (pvp *PVPFight) GetType() messageType {
	return pvp.Type
}

type Loot struct {
	Type messageType
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
	Type messageType
}

func (e *Error) GetType() messageType {
	return e.Type
}

type Accept struct {
	Type messageType
}

func (a *Accept) GetType() messageType {
	return a.Type
}

type Room struct {
	Type messageType
}

func (r *Room) GetType() messageType {
	return r.Type
}

type Character struct {
	Type messageType
}

func (c *Character) GetType() messageType {
	return c.Type
}

type Game struct {
	Type messageType
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
	Type messageType
}

func (c *Connection) GetType() messageType {
	return c.Type
}

type Version struct {
	Type messageType
}

func (v *Version) GetType() messageType {
	return v.Type
}
