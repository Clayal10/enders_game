package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/Clayal10/enders_game/lib/lurk"
)

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
		c.q <- lurkMessage
	}
}

// Will attempt to read from the message queue and will time out after 5 seconds.
func (c *Client) timeoutChannelRead() lurk.LurkMessage {
	for {
		select {
		case msg := <-c.q:
			return msg
		case <-time.After(time.Millisecond * 5000): // experiment with this
			return nil
		}
	}
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

	type jsonCharacter struct {
		Name        string `json:"name"`
		Attack      string `json:"attack"`
		Defense     string `json:"defense"`
		Regen       string `json:"regen"`
		Description string `json:"description"`
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
