package client

import (
	"errors"
	"net"
	"os"
	"sync"
	"time"

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

	mu              sync.Mutex
	conn            net.Conn
	initialMessages []lurk.LurkMessage
}

func Setup(cfg *ServerConfig) (*Client, error) {
	return connectToServer(cfg)
}

// Start should be called in a separate routine.
func (client *Client) Start() {

}

func connectToServer(cfg *ServerConfig) (*Client, error) {
	conn, err := net.Dial("tcp", cfg.Hostname+":"+cfg.Port)
	if err != nil {
		return nil, err
	}

	lurkMessages, err := func(conn net.Conn) (messages []lurk.LurkMessage, _ error) {
		for {
			conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
			ba, _, err := lurk.ReadSingleMessage(conn)
			if err != nil {
				if errors.Is(err, os.ErrDeadlineExceeded) {
					return messages, nil
				}
				return nil, err
			}

			lmsg, err := lurk.Unmarshal(ba)
			if err != nil {
				return nil, err
			}

			messages = append(messages, lmsg)
		}
	}(conn)
	if err != nil {
		return nil, err
	}

	gameMessage := func() lurk.LurkMessage {
		for _, lm := range lurkMessages {
			if lm.GetType() == lurk.TypeGame {
				return lm
			}
		}
		return nil
	}()

	if gameMessage == nil {
		return nil, errors.New("no game message")
	}

	game := gameMessage.(*lurk.Game)
	return &Client{
		Game:            game,
		initialMessages: lurkMessages,
		conn:            conn,
	}, nil
}
