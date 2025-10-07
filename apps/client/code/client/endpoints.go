package client

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

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
	})
}

const updateEP = "/lurk-client/update/"

func (c *Client) registerUpdateEP() {
	http.HandleFunc(fmt.Sprintf("%s%d/", updateEP, c.id), func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			return
		}

		messages := c.dequeueAll()
		if len(messages) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		c.updateClientState(messages)

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
			return
		}
		w.WriteHeader(http.StatusOK)
		log.Printf("ID:%v terminated from client", c.id)
		c.cf()
	})
}

const changeRoomEP = "/lurk-client/change-room/"

func (c *Client) registerChangeRoomEP() {
	http.HandleFunc(fmt.Sprintf("%s%d/", changeRoomEP, c.id), func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			return
		}

		type jsonChangeRoom struct {
			RoomNumber string `json:"roomNumber"`
		}

		ba, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		ch := &jsonChangeRoom{}
		json.Unmarshal(ba, ch)
		roomNum, err := strconv.ParseInt(ch.RoomNumber, 10, 16)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if _, err = c.conn.Write(lurk.Marshal(&lurk.ChangeRoom{
			RoomNumber: uint16(roomNum),
		})); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

}

const fightEP = "/lurk-client/fight/"
const pvpFightEP = "/lurk-client/pvp/"
const messageEP = "/lurk-client/message"
