package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/Clayal10/enders_game/cmd/server/code/server"
	"github.com/Clayal10/enders_game/pkg/assert"
	"github.com/Clayal10/enders_game/pkg/cross"
	"github.com/Clayal10/enders_game/pkg/lurk"
)

func TestHittingEndpoints(t *testing.T) {
	a := assert.New(t)
	serverPort := cross.GetFreePort()
	serverConfig := &server.Config{
		Port: serverPort,
	}
	cfs, err := server.New(serverConfig)
	a.NoError(err)
	pollTime = time.Millisecond
	defer func() {
		pollTime = 5 * time.Second
		for _, cf := range cfs {
			cf()
		}
	}()

	clientPort := cross.GetFreePort()
	//nolint:errcheck // In a test.
	go http.ListenAndServe(fmt.Sprintf("0.0.0.0:%v", clientPort), nil)

	client, err := New(&Config{
		Port: fmt.Sprint(serverPort),
	})
	a.NoError(err)
	client.Start()

	_ = startClientConnection(a, serverConfig, &lurk.Character{
		Name:       "Message Guy",
		Attack:     10,
		Defense:    20,
		Regen:      30,
		PlayerDesc: "A bot who is going to receive a message",
	})

	tests := createGameActions(fmt.Sprint(client.id))
	for _, test := range tests {
		t.Run(test.name, func(_ *testing.T) {
			resp, err := http.Post(fmt.Sprintf("http://localhost:%v%s", clientPort, test.endpoint), "application/json", bytes.NewBuffer(test.payload))
			a.NoError(err)
			a.True(resp.StatusCode == test.expected)
			_ = resp.Body.Close()
		})
	}
}

type gameAction struct {
	name     string
	endpoint string
	expected int
	payload  []byte
}

func createGameActions(id string) []gameAction {
	id += "/"
	startBA, _ := json.Marshal(&jsonCharacter{
		"tester", "25", "25", "25", "yes", "test",
	})
	startBadBA, _ := json.Marshal(&jsonCharacter{
		"tester", "huh", "25", "25", "", "test",
	})
	return []gameAction{
		{
			"happy start",
			startEP + id,
			http.StatusOK,
			startBA,
		},
		{
			"unhappy start",
			startEP + id,
			http.StatusBadRequest,
			startBadBA,
		},
		{
			"update happy",
			updateEP + id,
			http.StatusOK,
			[]byte("{}"),
		},
		{
			"update no content",
			updateEP + id,
			http.StatusNoContent,
			[]byte("{}"),
		},
		{
			"changeroom happy",
			changeRoomEP + id,
			http.StatusOK,
			[]byte(`{"roomNumber": "3"}`),
		},
		{
			"changeroom unhappy",
			changeRoomEP + id,
			http.StatusBadRequest,
			[]byte(`{"roomNumber": "huh"}`),
		},
		{
			"fight happy",
			fightEP + id,
			http.StatusOK,
			[]byte("{}"),
		},
		{
			"loot happy",
			lootEP + id,
			http.StatusOK,
			[]byte(`{"target": "test"}`),
		},
		{
			"loot unhappy",
			lootEP + id,
			http.StatusBadRequest,
			[]byte(`{"Targe"}`),
		},
		{
			"pvp happy",
			pvpFightEP + id,
			http.StatusOK,
			[]byte(`{"target": "test"}`),
		},
		{
			"pvp unhappy",
			pvpFightEP + id,
			http.StatusBadRequest,
			[]byte(`{"tar "test"`),
		},
		{
			"message happy",
			messageEP + id,
			http.StatusOK,
			[]byte(`{
				"recipient": "Message Guy",
				"text": "Test!"
			}`),
		},
		{
			"message unhappy",
			messageEP + id,
			http.StatusBadRequest,
			[]byte(`"Test!"
			}`),
		},
		{
			"update happy",
			updateEP + id,
			http.StatusOK,
			[]byte("{}"),
		},
		{
			"happy terminate", // call at the end
			terminateEP + id,
			http.StatusOK,
			[]byte("{}"),
		},
	}
}
