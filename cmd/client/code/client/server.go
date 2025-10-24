package client

import (
	"context"
	"log"
	"net"
	"time"
)

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

	c := newClient(conn, id)

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
	c.registerChangeRoomEP()
	c.registerFightEP()
	c.registerLootEP()
	c.registerPvpEP()
	c.registerMessageEP()
	// register more.
	log.Printf("Registered endpoints for ID:%d", c.id)
}
