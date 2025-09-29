package client

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Clayal10/enders_game/lib/lurk"
)

// Endpoints used for post methods will not write back any data to the UI.
// The response data will be used with the update endpoint.

const startEP = "/lurk-client/start/"

func (c *Client) registerStartEP() {
	http.HandleFunc(fmt.Sprintf("%s%d/", startEP, c.id), func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			return
		}

		char, err := c.getOrMakeCharacter(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		c.character = char
		c.character.Flags[lurk.Alive] = true

		if _, err = c.conn.Write(lurk.Marshal(c.character)); err != nil {
			log.Printf("%s: could not write Character to server", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, err = c.conn.Write(lurk.Marshal(&lurk.Start{}))
		if err != nil {
			log.Printf("%s: could not write start to server", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

const updateEP = "/lurk-client/update/"

func (c *Client) registerUpdateEP() {
	http.HandleFunc(fmt.Sprintf("%s%d/", updateEP, c.id), func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			return
		}

		msg := c.timeoutChannelRead()
		if msg == nil {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		c.updateClientState([]lurk.LurkMessage{msg})

		jsonString, err := json.Marshal(c.State)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, err = w.Write(jsonString)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}

const terminateEP = "/lurk-client/terminate/"

// This endpoint shall be called when the page is closed.
func (c *Client) registerTerminateEP() {
	http.HandleFunc(fmt.Sprintf("%s%d/", terminateEP, c.id), func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			return
		}

		_, err := c.conn.Write(lurk.Marshal(&lurk.Leave{}))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusOK)
		log.Printf("ID:%v terminated from client", c.id)
		c.cf()
	})
}
