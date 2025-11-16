package client

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/Clayal10/enders_game/pkg/lurk"
)

// Endpoints used for post methods will not write back any data to the UI.
// The response data will be used with the update endpoint.

const startEP = "/lurk-client/start/"

func (c *Client) registerStartEP() {
	http.HandleFunc(fmt.Sprintf("%s%d/", startEP, c.id), func(w http.ResponseWriter, r *http.Request) {
		char, err := c.getOrMakeCharacter(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
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

type jsonChangeRoom struct {
	RoomNumber string `json:"roomNumber"`
}

func (c *Client) registerChangeRoomEP() {
	http.HandleFunc(fmt.Sprintf("%s%d/", changeRoomEP, c.id), func(w http.ResponseWriter, r *http.Request) {
		ba, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		ch := &jsonChangeRoom{}
		if err := json.Unmarshal(ba, ch); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
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

func (c *Client) registerFightEP() {
	http.HandleFunc(fmt.Sprintf("%s%d/", fightEP, c.id), func(w http.ResponseWriter, r *http.Request) {
		if _, err := c.conn.Write(lurk.Marshal(&lurk.Fight{})); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
}

const lootEP = "/lurk-client/loot/"

type jsonLoot struct {
	TargetName string `json:"target"`
}

func (c *Client) registerLootEP() {
	http.HandleFunc(fmt.Sprintf("%s%d/", lootEP, c.id), func(w http.ResponseWriter, r *http.Request) {
		ba, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		loot := &jsonLoot{}
		if err = json.Unmarshal(ba, loot); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if _, err := c.conn.Write(lurk.Marshal(&lurk.Loot{
			TargetName: loot.TargetName,
		})); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
}

const pvpFightEP = "/lurk-client/pvp/"

type jsonPVP struct {
	TargetName string `json:"target"`
}

func (c *Client) registerPvpEP() {
	http.HandleFunc(fmt.Sprintf("%s%d/", pvpFightEP, c.id), func(w http.ResponseWriter, r *http.Request) {
		ba, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		pvp := &jsonPVP{}
		if err = json.Unmarshal(ba, pvp); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if _, err := c.conn.Write(lurk.Marshal(&lurk.PVPFight{
			TargetName: pvp.TargetName,
		})); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
}

const messageEP = "/lurk-client/message/"

type jsonMessage struct {
	Recipient string `json:"recipient"`
	Text      string `json:"text"`
}

func (c *Client) registerMessageEP() {
	http.HandleFunc(fmt.Sprintf("%s%d/", messageEP, c.id), func(w http.ResponseWriter, r *http.Request) {
		ba, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		msg := &jsonMessage{}
		if err = json.Unmarshal(ba, msg); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if _, err := c.conn.Write(lurk.Marshal(&lurk.Message{
			Recipient: msg.Recipient,
			Sender:    c.character.Name,
			Text:      msg.Text,
		})); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}
