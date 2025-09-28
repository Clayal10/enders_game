package client

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/Clayal10/enders_game/lib/lurk"
)

// Each field correlates to a section of the UI that should be updated in the JS.
// This struct will be json'd and each field should be added to the innerHTML of the HTML element,
// and should not overwrite the data already in it.
type ClientState struct {
	Info    string `json:"info"`
	Rooms   string `json:"rooms"`
	Players string `json:"players"`
	Id      int64  `json:"id"`

	characters map[string]*lurk.Character
	rooms      map[uint16]*lurk.Room
	// connections too.
}

// Client contains all data needed to run a client instance.
type Client struct {
	Game      *lurk.Game
	character *lurk.Character
	State     *ClientState

	id  int64
	ctx context.Context
	cf  context.CancelFunc

	//mu   sync.Mutex
	conn net.Conn
	q    chan lurk.LurkMessage
}

type Config struct {
	Hostname, Port string
}

func New(cfg *Config) (*Client, error) {
	conn, err := net.Dial("tcp", cfg.Hostname+":"+cfg.Port)
	if err != nil {
		return nil, err
	}

	lurkMessages, err := readAllMessagesInBuffer(conn)
	if err != nil {
		return nil, err
	}

	id := time.Now().UnixMicro()

	c := &Client{
		conn: conn,
		id:   id,
		q:    make(chan lurk.LurkMessage, 100),
		State: &ClientState{
			Id:         id,
			characters: map[string]*lurk.Character{},
			rooms:      map[uint16]*lurk.Room{},
		},
	}

	c.updateClientState(lurkMessages)

	return c, nil
}

func (c *Client) Start() {
	c.registerEndpoints()
	c.ctx, c.cf = context.WithCancel(context.Background())
	go c.readFromServer()
}

func (c *Client) registerEndpoints() {
	c.registerStartEP()
	c.registerUpdateEP()
	c.registerTerminateEP()
	// register more.
	log.Printf("Registered endpoints for ID:%d", c.id)
}
