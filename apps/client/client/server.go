package client

import (
	"sync"

	"github.com/Clayal10/enders_game/lib/lurk"
)

type ServerConfig struct {
	Hostname string
	Port     string
}

// Client contains all data needed to run a client instance.
type Client struct {
	// Message queue?

	Game *lurk.Game

	mu sync.Mutex
}

func Setup(cfg *ServerConfig) (*Client, error) {
	if err := validateConfig(cfg); err != nil {
		return nil, err
	}
	return &Client{
		Game: &lurk.Game{
			Type:          lurk.TypeGame,
			InitialPoints: 100,
			StatLimit:     65000,
			GameDesc:      "Placeholder Test!",
		},
	}, nil
}

// Start should be called in a separate routine.
func (client *Client) Start() {

}

func validateConfig(cfg *ServerConfig) error {
	return nil
}
