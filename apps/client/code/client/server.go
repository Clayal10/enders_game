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
type ClientUpdate struct {
	Info    string `json:"info"`
	Rooms   string `json:"rooms"`
	Players string `json:"players"`
	Id      int64  `json:"id"`
}

// Client contains all data needed to run a client instance.
type Client struct {
	Game      *lurk.Game
	character *lurk.Character

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

func New(cfg *Config) (*Client, *ClientUpdate, error) {
	conn, err := net.Dial("tcp", cfg.Hostname+":"+cfg.Port)
	if err != nil {
		return nil, nil, err
	}

	lurkMessages, err := readAllMessagesInBuffer(conn)
	if err != nil {
		return nil, nil, err
	}

	c := &Client{
		conn: conn,
		id:   time.Now().UnixMicro(),
		q:    make(chan lurk.LurkMessage, 100),
	}

	cu := c.getClientUpdate(lurkMessages)

	return c, cu, nil
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
