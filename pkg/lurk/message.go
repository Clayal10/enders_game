package lurk

import (
	"encoding/binary"
	"net"
	"time"

	"github.com/Clayal10/enders_game/lib/cross"
)

type MessageType byte

// Exported message types.
const (
	TypeMessage    MessageType = 1
	TypeChangeRoom MessageType = 2
	TypeFight      MessageType = 3
	TypePVPFight   MessageType = 4
	TypeLoot       MessageType = 5
	TypeStart      MessageType = 6
	TypeError      MessageType = 7
	TypeAccept     MessageType = 8
	TypeRoom       MessageType = 9
	TypeCharacter  MessageType = 10
	TypeGame       MessageType = 11
	TypeLeave      MessageType = 12
	TypeConnection MessageType = 13
	TypeVersion    MessageType = 14
)

// LengthOffset is a key that will tell you how many bytes you will need to read per message
// type to have a full enough message. Fields not denoted with '// X' have fixed length messages
// and the returned value is good. Otherwise, send it through the 'GetVariableRate' function
// and find out the real total.
//
// All length fields are 16 bits in each message.
var LengthOffset = map[MessageType]int{
	TypeMessage:    3,
	TypeChangeRoom: 3,  // X
	TypeFight:      1,  // X
	TypePVPFight:   33, // X
	TypeLoot:       33, // X
	TypeStart:      1,  // X
	TypeError:      4,
	TypeAccept:     2, // X
	TypeRoom:       37,
	TypeCharacter:  48,
	TypeGame:       7,
	TypeLeave:      1, // X
	TypeConnection: 37,
	TypeVersion:    5,
}

const (
	maxStringLen = 32
	// length of variable length messages  before their text.
	messageLength = 67
)

// Exported character flags
const (
	Alive      = "Alive"
	JoinBattle = "Join Battle"
	Monster    = "Monster"
	Started    = "Started"
	Ready      = "Ready"
)

// Those using this library will need to use this function on the returned interface
// to know which type will need to be used for assertion.
type LurkMessage interface {
	GetType() MessageType
}

// GetVariableLength will return the total byte length of the message based off of
// the variable length message.
func GetVariableLength(data []byte) (int, error) {
	if err := validate(data); err != nil {
		return 0, err
	}

	msgType := MessageType(data[0])

	idx, ok := LengthOffset[msgType]
	if !ok {
		return 0, cross.ErrInvalidMessageType
	}

	if len(data) < idx {
		return 0, cross.ErrFrameTooSmall
	}

	switch msgType {
	case TypeMessage:
		msgLength := int(binary.LittleEndian.Uint16(data[idx-2:]))
		return msgLength, nil
	case TypeError:
		msgLength := int(binary.LittleEndian.Uint16(data[idx-2:]))
		return msgLength, nil
	case TypeRoom:
		msgLength := int(binary.LittleEndian.Uint16(data[idx-2:]))
		return msgLength, nil
	case TypeCharacter:
		msgLength := int(binary.LittleEndian.Uint16(data[idx-2:]))
		return msgLength, nil
	case TypeConnection:
		msgLength := int(binary.LittleEndian.Uint16(data[idx-2:]))
		return msgLength, nil
	case TypeVersion:
		msgLength := int(binary.LittleEndian.Uint16(data[idx-2:]))
		return msgLength, nil
	case TypeGame:
		msgLength := int(binary.LittleEndian.Uint16(data[idx-2:]))
		return msgLength, nil
	default:
		return -1, nil
	}
}

// Unmarshal takes a slice of bytes and returns a LurkMessage interface object. The
// LurkMessage then needs to be type asserted based on the type returned from the
// 'GetType()' function.
func Unmarshal(data []byte) (LurkMessage, error) {
	if err := validate(data); err != nil {
		return nil, err
	}
	// Various unmarshaling and returning of their respective types.
	switch MessageType(data[0]) {
	case TypeMessage:
		return unmarshalMessage(data)
	case TypeChangeRoom:
		return unmarshalChangeRoom(data)
	case TypeFight:
		return &Fight{Type: TypeFight}, nil
	case TypePVPFight:
		return unmarshalPVP(data)
	case TypeLoot:
		return unmarshalLoot(data)
	case TypeStart:
		return &Start{Type: TypeStart}, nil
	case TypeError:
		return unmarshalError(data)
	case TypeAccept:
		return unmarshalAccept(data)
	case TypeRoom:
		return unmarshalRoom(data)
	case TypeCharacter:
		return unmarshalCharacter(data)
	case TypeGame:
		return unmarshalGame(data)
	case TypeLeave:
		return &Leave{Type: TypeLeave}, nil
	case TypeConnection:
		return unmarshalConnection(data)
	case TypeVersion:
		return unmarshalVersion(data)
	}
	return nil, cross.ErrInvalidMessageType
}

// Marshal Will take any LurkMessage object and return a byte array
// ready for messaging.
func Marshal(lm LurkMessage) []byte {
	switch lm.GetType() {
	case TypeMessage:
		if msg, ok := lm.(*Message); ok {
			return marshalMessage(msg)
		}
	case TypeChangeRoom:
		if cr, ok := lm.(*ChangeRoom); ok {
			return marshalChangeRoom(cr)
		}
	case TypeFight:
		return []byte{0x03}
	case TypePVPFight:
		if pvp, ok := lm.(*PVPFight); ok {
			return marshalPVP(pvp)
		}
	case TypeLoot:
		if l, ok := lm.(*Loot); ok {
			return marshalLoot(l)
		}
	case TypeStart:
		return []byte{0x06}
	case TypeError:
		if e, ok := lm.(*Error); ok {
			return marshalError(e)
		}
	case TypeAccept:
		if a, ok := lm.(*Accept); ok {
			return marshalAccept(a)
		}
	case TypeRoom:
		if room, ok := lm.(*Room); ok {
			return marshalRoom(room)
		}
	case TypeCharacter:
		if char, ok := lm.(*Character); ok {
			return marshalCharacter(char)
		}
	case TypeGame:
		if game, ok := lm.(*Game); ok {
			return marshalGame(game)
		}
	case TypeLeave:
		return []byte{0xc}
	case TypeConnection:
		if conn, ok := lm.(*Connection); ok {
			return marshalConnection(conn)
		}
	case TypeVersion:
		if e, ok := lm.(*Version); ok {
			return marshalVersion(e)
		}
	}
	return nil
}

// We want to read exactly the length of the message. This function will do up to 3
// calls to 'Read' to read exactly one message.
func ReadSingleMessage(conn net.Conn) ([]byte, int, error) {
	const messageTypePadding = 64
	buffer := make([]byte, 1)
	if _, err := conn.Read(buffer); err != nil {
		return nil, 0, err
	}

	bytesNeeded, ok := LengthOffset[MessageType(buffer[0])]
	if !ok {
		return nil, 0, cross.ErrInvalidMessageType
	}
	if MessageType(buffer[0]) == TypeMessage {
		bytesNeeded += messageTypePadding
	}

	if bytesNeeded == 1 {
		return buffer, 1, nil
	}

	b := make([]byte, bytesNeeded-1)
	n := 0
	defer func() { _ = conn.SetReadDeadline(time.Time{}) }()
	for n < len(b) {
		_ = conn.SetReadDeadline(time.Now().Add(1000 * time.Millisecond))
		m, err := conn.Read(b)
		if err != nil {
			return nil, 0, err
		}
		buffer = append(buffer, b[:m]...)
		n += m
	}

	varLen, err := GetVariableLength(buffer)
	if err != nil {
		return nil, 0, err
	}

	if varLen == -1 {
		return buffer, bytesNeeded, nil
	}

	b = make([]byte, varLen)
	n = 0
	for n < len(b) {
		_ = conn.SetReadDeadline(time.Now().Add(1000 * time.Millisecond))
		m, err := conn.Read(b)
		if err != nil {
			return nil, 0, err
		}
		buffer = append(buffer, b[:m]...)
		n += m
	}

	return buffer, varLen + bytesNeeded, nil
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
	Type      MessageType
	Recipient string // max 32 bytes. All fields noted with bytes are null terminated '\x00'.
	Sender    string // max 30 bytes
	Text      string
	Narration bool
}

func (*Message) GetType() MessageType {
	return TypeMessage
}

func unmarshalMessage(data []byte) (*Message, error) {
	if len(data) < 67 {
		return nil, cross.ErrFrameTooSmall
	}
	m := &Message{}
	offset := 0
	m.Type = MessageType(data[offset])
	offset++
	msgLen := binary.LittleEndian.Uint16(data[offset:])
	offset += 2

	if len(data) < 67+int(msgLen) {
		return nil, cross.ErrFrameTooSmall
	}

	l := getNullTermLen(data[offset:])
	m.Recipient = string(data[offset : offset+l])
	offset += maxStringLen
	l = getNullTermLen(data[offset:])
	m.Sender = string(data[offset : offset+l])
	offset += maxStringLen - 1

	// Check for narration
	if data[offset] == 1 {
		m.Narration = true
	}
	offset++
	m.Text = string(data[offset : offset+int(msgLen)])

	return m, nil
}

func marshalMessage(msg *Message) []byte {
	msgLength := uint16(len(msg.Text))
	ba := make([]byte, msgLength+messageLength)

	offset := 0
	ba[offset] = byte(TypeMessage)
	offset++
	binary.LittleEndian.PutUint16(ba[offset:], msgLength)
	offset += 2
	copy(ba[offset:offset+maxStringLen], getNullTermedString(msg.Recipient))
	offset += maxStringLen
	copy(ba[offset:offset+maxStringLen], getNullTermedString(msg.Sender))
	offset += maxStringLen - 1
	ba[offset] = boolToByte(msg.Narration)
	offset++
	copy(ba[offset:], []byte(msg.Text))
	return ba
}

type ChangeRoom struct {
	Type       MessageType
	RoomNumber uint16
}

func (*ChangeRoom) GetType() MessageType {
	return TypeChangeRoom
}

func unmarshalChangeRoom(data []byte) (*ChangeRoom, error) {
	if len(data) < 3 {
		return nil, cross.ErrFrameTooSmall
	}
	return &ChangeRoom{
		Type:       MessageType(data[0]),
		RoomNumber: binary.LittleEndian.Uint16(data[1:]),
	}, nil
}

func marshalChangeRoom(cr *ChangeRoom) []byte {
	ba := make([]byte, 3)
	ba[0] = byte(TypeChangeRoom)
	binary.LittleEndian.PutUint16(ba[1:], cr.RoomNumber)
	return ba
}

type Fight struct {
	Type MessageType
}

func (*Fight) GetType() MessageType {
	return TypeFight
}

type PVPFight struct {
	Type       MessageType
	TargetName string // 32 bytes
}

func (*PVPFight) GetType() MessageType {
	return TypePVPFight
}

func unmarshalPVP(data []byte) (*PVPFight, error) {
	if len(data) < maxStringLen+1 {
		return nil, cross.ErrFrameTooSmall
	}

	nameLen := getNullTermLen(data[1:])

	return &PVPFight{
		Type:       MessageType(data[0]),
		TargetName: string(data[1 : nameLen+1]),
	}, nil
}

func marshalPVP(pvp *PVPFight) []byte {
	ba := make([]byte, maxStringLen+1)
	ba[0] = byte(TypePVPFight)
	copy(ba[1:], getNullTermedString(pvp.TargetName))
	return ba
}

type Loot struct {
	Type       MessageType
	TargetName string // 32 bytes
}

func (l *Loot) GetType() MessageType {
	return TypeLoot
}

func unmarshalLoot(data []byte) (*Loot, error) {
	if len(data) < maxStringLen+1 {
		return nil, cross.ErrFrameTooSmall
	}

	nameLen := getNullTermLen(data[1:])

	return &Loot{
		Type:       MessageType(data[0]),
		TargetName: string(data[1 : nameLen+1]),
	}, nil
}

func marshalLoot(loot *Loot) []byte {
	ba := make([]byte, maxStringLen+1)
	ba[0] = byte(TypeLoot)
	copy(ba[1:], getNullTermedString(loot.TargetName))
	return ba
}

type Start struct {
	Type MessageType
}

func (s *Start) GetType() MessageType {
	return TypeStart
}

type Error struct {
	Type       MessageType
	ErrCode    cross.ErrCode
	ErrMessage string
}

func (e *Error) GetType() MessageType {
	return TypeError
}

func unmarshalError(data []byte) (*Error, error) {
	if len(data) < 4 {
		return nil, cross.ErrFrameTooSmall
	}

	msgLen := binary.LittleEndian.Uint16(data[2:])

	if len(data) < int(4+msgLen) {
		return nil, cross.ErrFrameTooSmall
	}

	errCode := cross.ErrCode(data[1])
	if errCode > cross.NoPVP {
		return nil, cross.ErrInvalidErrCode
	}

	return &Error{
		Type:       MessageType(data[0]),
		ErrCode:    errCode,
		ErrMessage: string(data[4 : 4+msgLen]),
	}, nil
}

func marshalError(e *Error) []byte {
	ba := make([]byte, 4+len(e.ErrMessage))
	ba[0] = byte(TypeError)
	ba[1] = byte(e.ErrCode)
	binary.LittleEndian.PutUint16(ba[2:], uint16(len(e.ErrMessage)))
	copy(ba[4:], []byte(e.ErrMessage))
	return ba
}

type Accept struct {
	Type   MessageType
	Action MessageType
}

func (a *Accept) GetType() MessageType {
	return TypeAccept
}

func unmarshalAccept(data []byte) (*Accept, error) {
	if len(data) < 2 {
		return nil, cross.ErrFrameTooSmall
	}
	return &Accept{
		Type:   MessageType(data[0]),
		Action: MessageType(data[1]),
	}, nil
}

func marshalAccept(a *Accept) []byte {
	ba := make([]byte, 2)
	ba[0] = byte(TypeAccept)
	ba[1] = byte(a.Action)
	return ba
}

type Room struct {
	Type       MessageType
	RoomNumber uint16
	RoomName   string // 32 bytes
	RoomDesc   string
}

func (r *Room) GetType() MessageType {
	return TypeRoom
}

func unmarshalRoom(data []byte) (*Room, error) {
	if len(data) < 37 {
		return nil, cross.ErrFrameTooSmall
	}

	nameLen := getNullTermLen(data[3:])

	room := &Room{
		Type:       MessageType(data[0]),
		RoomNumber: binary.LittleEndian.Uint16(data[1:]),
		RoomName:   string(data[3 : 3+nameLen]),
	}

	offset := 3 + maxStringLen
	descLen := binary.LittleEndian.Uint16(data[offset:])
	if len(data) < int(37+descLen) {
		return nil, cross.ErrFrameTooSmall
	}

	offset += 2
	room.RoomDesc = string(data[offset : offset+int(descLen)])
	return room, nil
}

func marshalRoom(room *Room) []byte {
	ba := make([]byte, 37+len(room.RoomDesc))
	offset := 0
	ba[offset] = byte(TypeRoom)
	offset++
	binary.LittleEndian.PutUint16(ba[offset:], room.RoomNumber)
	offset += 2
	copy(ba[offset:], getNullTermedString(room.RoomName))
	offset += maxStringLen
	binary.LittleEndian.PutUint16(ba[offset:], uint16(len(room.RoomDesc)))
	offset += 2
	copy(ba[offset:], []byte(room.RoomDesc))
	return ba
}

type Character struct {
	Type       MessageType
	Name       string          // 32 bytes
	Flags      map[string]bool // Alive, Join, Monster, Started, Ready
	Attack     uint16
	Defense    uint16
	Regen      uint16
	Health     int16
	Gold       uint16
	RoomNum    uint16
	PlayerDesc string
}

func (c *Character) GetType() MessageType {
	return TypeCharacter
}

func unmarshalCharacter(data []byte) (*Character, error) {
	if len(data) < 48 {
		return nil, cross.ErrFrameTooSmall
	}

	c := &Character{
		Type: MessageType(data[0]),
	}

	nameLen := getNullTermLen(data[1:])
	c.Name = string(data[1 : nameLen+1])

	offset := 1 + maxStringLen
	c.Flags = unmarshalCharacterFlags(data[offset])
	offset++
	c.Attack = binary.LittleEndian.Uint16(data[offset:])
	offset += 2
	c.Defense = binary.LittleEndian.Uint16(data[offset:])
	offset += 2
	c.Regen = binary.LittleEndian.Uint16(data[offset:])
	offset += 2
	c.Health = int16(binary.LittleEndian.Uint16(data[offset:]))
	offset += 2
	c.Gold = binary.LittleEndian.Uint16(data[offset:])
	offset += 2
	c.RoomNum = binary.LittleEndian.Uint16(data[offset:])
	offset += 2
	descLen := binary.LittleEndian.Uint16(data[offset:])
	offset += 2
	if len(data) < int(48+descLen) {
		return nil, cross.ErrFrameTooSmall
	}
	c.PlayerDesc = string(data[offset : offset+int(descLen)])
	return c, nil
}

func unmarshalCharacterFlags(data byte) map[string]bool {
	flags := make(map[string]bool)
	flags[Alive] = aliveBit&data == aliveBit
	flags[JoinBattle] = joinBit&data == joinBit
	flags[Monster] = monsterBit&data == monsterBit
	flags[Started] = startedBit&data == startedBit
	flags[Ready] = readyBit&data == readyBit
	return flags
}

func marshalCharacter(c *Character) []byte {
	ba := make([]byte, 48+len(c.PlayerDesc))
	offset := 0
	ba[offset] = byte(TypeCharacter)
	offset++
	copy(ba[offset:], getNullTermedString(c.Name))
	offset += maxStringLen
	ba[offset] = marshalCharacterFlags(c.Flags)
	offset++
	binary.LittleEndian.PutUint16(ba[offset:], c.Attack)
	offset += 2
	binary.LittleEndian.PutUint16(ba[offset:], c.Defense)
	offset += 2
	binary.LittleEndian.PutUint16(ba[offset:], c.Regen)
	offset += 2
	binary.LittleEndian.PutUint16(ba[offset:], uint16(c.Health))
	offset += 2
	binary.LittleEndian.PutUint16(ba[offset:], c.Gold)
	offset += 2
	binary.LittleEndian.PutUint16(ba[offset:], c.RoomNum)
	offset += 2
	binary.LittleEndian.PutUint16(ba[offset:], uint16(len(c.PlayerDesc)))
	offset += 2
	copy(ba[offset:], []byte(c.PlayerDesc))
	return ba
}

func marshalCharacterFlags(flags map[string]bool) (word byte) {
	if flags[Alive] {
		word += aliveBit
	}
	if flags[JoinBattle] {
		word += joinBit
	}
	if flags[Monster] {
		word += monsterBit
	}
	if flags[Started] {
		word += startedBit
	}
	if flags[Ready] {
		word += readyBit
	}
	return
}

const (
	aliveBit   = 128
	joinBit    = 64
	monsterBit = 32
	startedBit = 16
	readyBit   = 8
)

type Game struct {
	Type          MessageType
	InitialPoints uint16
	StatLimit     uint16
	GameDesc      string
}

func (g *Game) GetType() MessageType {
	return TypeGame
}

func unmarshalGame(data []byte) (*Game, error) {
	if len(data) < 7 {
		return nil, cross.ErrFrameTooSmall
	}
	g := &Game{
		Type: MessageType(data[0]),
	}
	offset := 1
	g.InitialPoints = binary.LittleEndian.Uint16(data[offset:])
	offset += 2
	g.StatLimit = binary.LittleEndian.Uint16(data[offset:])
	offset += 2
	descLen := binary.LittleEndian.Uint16(data[offset:])
	offset += 2

	if len(data) < int(7+descLen) {
		return nil, cross.ErrFrameTooSmall
	}

	g.GameDesc = string(data[offset : offset+int(descLen)])
	return g, nil
}

func marshalGame(g *Game) []byte {
	ba := make([]byte, 7+len(g.GameDesc))
	offset := 0
	ba[offset] = byte(TypeGame)
	offset++
	binary.LittleEndian.PutUint16(ba[offset:], g.InitialPoints)
	offset += 2
	binary.LittleEndian.PutUint16(ba[offset:], g.StatLimit)
	offset += 2
	binary.LittleEndian.PutUint16(ba[offset:], uint16(len(g.GameDesc)))
	offset += 2
	copy(ba[offset:], []byte(g.GameDesc))
	return ba
}

type Leave struct {
	Type MessageType
}

func (l *Leave) GetType() MessageType {
	return TypeLeave
}

type Connection struct {
	Type       MessageType
	RoomNumber uint16
	RoomName   string //32 bytes
	RoomDesc   string
}

func (c *Connection) GetType() MessageType {
	return TypeConnection
}

func unmarshalConnection(data []byte) (*Connection, error) {
	if len(data) < 37 {
		return nil, cross.ErrFrameTooSmall
	}
	c := &Connection{
		Type:       MessageType(data[0]),
		RoomNumber: binary.LittleEndian.Uint16(data[1:]),
	}
	offset := 3
	nameLen := getNullTermLen(data[offset:])
	c.RoomName = string(data[offset : offset+nameLen])
	offset += maxStringLen
	descLen := binary.LittleEndian.Uint16(data[offset:])
	if len(data) < int(37+descLen) {
		return nil, cross.ErrFrameTooSmall
	}
	offset += 2
	c.RoomDesc = string(data[offset : offset+int(descLen)])
	return c, nil
}

func marshalConnection(c *Connection) []byte {
	ba := make([]byte, 37+len(c.RoomDesc))
	offset := 0
	ba[offset] = byte(TypeConnection)
	offset++
	binary.LittleEndian.PutUint16(ba[offset:], c.RoomNumber)
	offset += 2
	copy(ba[offset:], getNullTermedString(c.RoomName))
	offset += maxStringLen
	binary.LittleEndian.PutUint16(ba[offset:], uint16(len(c.RoomDesc)))
	offset += 2
	copy(ba[offset:], []byte(c.RoomDesc))
	return ba
}

type Version struct {
	Type       MessageType
	Major      byte
	Minor      byte
	Extensions [][]byte // For now. Turn into object when we know what it is.
}

func (v *Version) GetType() MessageType {
	return TypeVersion
}

func unmarshalVersion(data []byte) (*Version, error) {
	if len(data) < 5 {
		return nil, cross.ErrFrameTooSmall
	}
	v := &Version{
		Type:  MessageType(data[0]),
		Major: data[1],
		Minor: data[2],
	}
	offset := 3
	listLen := binary.LittleEndian.Uint16(data[offset:])
	offset += 2
	if len(data) < offset+int(listLen) {
		return nil, cross.ErrFrameTooSmall
	}
	if listLen == 0 {
		return v, nil
	}
	for {
		if len(data) <= offset+1 {
			return v, nil
		}
		extLen := binary.LittleEndian.Uint16(data[offset:])
		offset += 2
		if len(data) < offset+int(extLen) {
			return v, nil
		}
		v.Extensions = append(v.Extensions, data[offset:offset+int(extLen)])
		offset += int(extLen)
	}
}

func marshalVersion(v *Version) []byte {
	ba := make([]byte, 5) // Only big enough for the size of the list of extensions
	offset := 0
	ba[offset] = byte(TypeVersion)
	offset++
	ba[offset] = v.Major
	offset++
	ba[offset] = v.Minor
	offset++

	remainingLen := 0
	for _, ext := range v.Extensions {
		remainingLen += len(ext) + 2 // for the size
	}

	binary.LittleEndian.PutUint16(ba[offset:], uint16(remainingLen))

	for _, ext := range v.Extensions {
		ba = binary.LittleEndian.AppendUint16(ba, uint16(len(ext)))
		ba = append(ba, ext...)
	}
	return ba
}

// data should be a slice starting at the start of a null terminated string.
func getNullTermLen(data []byte) (length int) {
	for _, b := range data {
		if b == 0 || length == maxStringLen {
			break
		}
		length++
	}
	return
}

func getNullTermedString(value string) []byte {
	length := len(value)
	if length >= maxStringLen {
		value = value[:maxStringLen]
	}

	nulls := make([]byte, maxStringLen-len(value))
	for i := range nulls {
		nulls[i] = 0
	}
	return append([]byte(value), nulls...)
}

func boolToByte(cond bool) byte {
	if cond {
		return 1
	}
	return 0
}
