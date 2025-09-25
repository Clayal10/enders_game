package client

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"sync"
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
	// Message queue?

	Game *lurk.Game
	id   int64

	mu   sync.Mutex
	conn net.Conn
}

type Config struct {
	Hostname, Port string
}

const startEP = "/lurk-client/start/"
const updateEP = "/lurk-client/update/"

func New(cfg *Config) (*Client, *ClientUpdate, error) {
	conn, err := net.Dial("tcp", cfg.Hostname+":"+cfg.Port)
	if err != nil {
		return nil, nil, err
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
		return nil, nil, err
	}

	c := &Client{
		conn: conn,
		id:   time.Now().UnixNano(),
	}

	cu := c.getClientUpdate(lurkMessages)

	return c, cu, nil
}

func (c *Client) Start() error {
	c.registerEndpoints()
	return nil
}

func (client *Client) getClientUpdate(lurkMessages []lurk.LurkMessage) *ClientUpdate {
	cu := &ClientUpdate{
		Id: client.id,
	}
	for _, msg := range lurkMessages {
		switch msg.GetType() {
		case lurk.TypeGame:
			game := msg.(*lurk.Game)
			client.Game = game
			cu.Info += fmt.Sprintf("Stat Limit: %v\nInitial Points: %v\n%s\n", game.StatLimit, game.InitialPoints, game.GameDesc)
		case lurk.TypeVersion:
			version := msg.(*lurk.Version)
			extensionBytes := 0
			for _, v := range version.Extensions {
				extensionBytes += len(v)
			}
			cu.Info += fmt.Sprintf("LURK Version %v.%v | %v Bytes of Extensions\n", version.Major, version.Minor, extensionBytes)
		case lurk.TypeCharacter:
			character := msg.(*lurk.Character)
			cu.Players += fmt.Sprintf("%s | Attack: %v Defense: %v Health: %v Monster?: %v\n", character.Name, character.Attack, character.Defense, character.Health, character.Flags[lurk.Monster])
		}
	}
	return cu
}

func (c *Client) registerEndpoints() {
	http.HandleFunc(fmt.Sprintf("%s%d/", startEP, c.id), handleStart)
	log.Printf("Registered endpoints for ID:%d", c.id)
}

func handleStart(w http.ResponseWriter, r *http.Request) {}
