package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/Clayal10/enders_game/pkg/data"
	"github.com/Clayal10/enders_game/pkg/lurk"
)

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
	q    *data.Queue[lurk.LurkMessage]
}

func newClient(conn net.Conn, id int64) *Client {
	return &Client{
		conn:  conn,
		id:    id,
		q:     data.NewQueue[lurk.LurkMessage](100),
		State: newClientState(id),
	}
}

// This function enqueues messages into the message queue to be returned
// as HTTP responses to the UI.
func (c *Client) readFromServer() {
	var ba []byte
	var lurkMessage lurk.LurkMessage
	var err error
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
		}
		if ba, _, err = lurk.ReadSingleMessage(c.conn); err != nil {
			break
		}

		if lurkMessage, err = lurk.Unmarshal(ba); err != nil {
			break
		}
		c.q.Enqueue(lurkMessage)
	}
}

var pollTime = 5 * time.Second

func (c *Client) dequeueAll() (lms []lurk.LurkMessage) {
	start := time.Now()
	for time.Since(start) < pollTime {
		lm, err := c.q.Dequeue()
		if err != nil {
			if len(lms) != 0 {
				break
			}
			time.Sleep(2 * time.Millisecond) // In hopes of combining messages in a single http response.
			continue
		}
		lms = append(lms, lm)
	}
	return
}

func readAllMessagesInBuffer(conn net.Conn) (messages []lurk.LurkMessage, _ error) {
	for {
		conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
		defer func() {
			_ = conn.SetReadDeadline(time.Time{})
		}()
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
}

type jsonCharacter struct {
	Name        string `json:"name"`
	Attack      string `json:"attack"`
	Defense     string `json:"defense"`
	Regen       string `json:"regen"`
	Description string `json:"description"`
}

func (c *Client) getOrMakeCharacter(body io.ReadCloser) (*lurk.Character, error) {
	if body == nil {
		return &lurk.Character{
			Type: lurk.TypeCharacter,
			Name: fmt.Sprintf("Client %d", c.id),
			Flags: map[string]bool{
				lurk.Ready: true,
			},
			PlayerDesc: "Character",
		}, nil
	}

	ba, err := io.ReadAll(body)
	if err != nil {
		log.Printf("%v: could not read start body", err)
		return nil, err
	}
	jsonChar := &jsonCharacter{}
	if err = json.Unmarshal(ba, jsonChar); err != nil {
		log.Printf("%v: could not unmarshal into a LURK character", err)
		return nil, err
	}

	attack, err := strconv.ParseUint(jsonChar.Attack, 10, 16)
	if err != nil {
		return nil, err
	}
	defense, err := strconv.ParseUint(jsonChar.Defense, 10, 16)
	if err != nil {
		return nil, err
	}
	regen, err := strconv.ParseUint(jsonChar.Regen, 10, 16)
	if err != nil {
		return nil, err
	}
	return &lurk.Character{
		Name:       jsonChar.Name,
		Attack:     uint16(attack),
		Defense:    uint16(defense),
		Regen:      uint16(regen),
		PlayerDesc: jsonChar.Description,
		Flags: map[string]bool{
			lurk.Ready: true,
		},
	}, nil
}
