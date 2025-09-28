package client

import (
	"errors"
	"log"
	"net"
	"os"
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
		log.Printf("Got type %v from server", lurkMessage.GetType())
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
