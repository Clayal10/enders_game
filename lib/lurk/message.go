package lurk

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
}

// Unmarshal will take a raw frame and turn it into the appropriate type that
// satisfies the LurkMessage interface.
func Unmarshal(data []byte) LurkMessage

// Marshal Will take any LurkMessage object and return a byte array
// ready for messaging.
func Marshal(lm LurkMessage) []byte

type Message struct{}
type ChangeRoom struct{}
type Fight struct{}
type PVPFight struct{}
type Loot struct{}
type Start struct{}
type Error struct{}
type Accept struct{}
type Room struct{}
type Character struct{}
type Game struct{}
type Leave struct{}
type Connection struct{}
type Version struct{}
