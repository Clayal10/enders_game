package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/Clayal10/enders_game/cmd/server/code/server"
	"github.com/Clayal10/enders_game/pkg/assert"
	"github.com/Clayal10/enders_game/pkg/cross"
)

func TestHittingEndpoints(t *testing.T) {
	a := assert.New(t)
	serverPort := cross.GetFreePort()
	cfs, err := server.New(&server.Config{
		Port: serverPort,
	})
	a.NoError(err)
	defer func() {
		for _, cf := range cfs {
			cf()
		}
	}()

	clientPort := cross.GetFreePort()
	go http.ListenAndServe(fmt.Sprintf("0.0.0.0:%v", clientPort), nil)

	client, err := New(&Config{
		Port: fmt.Sprint(serverPort),
	})
	a.NoError(err)
	client.Start()
	tests := createTests(fmt.Sprint(client.id))
	for _, test := range tests {
		t.Run(test.name, func(_ *testing.T) {
			resp, err := http.Post(fmt.Sprintf("http://localhost:%v%s", clientPort, test.endpoint), "application/json", bytes.NewBuffer(test.payload))
			a.NoError(err)
			a.True(resp.StatusCode == test.expected)
			resp.Body.Close()
		})
	}
}

type test struct {
	name     string
	endpoint string
	expected int
	payload  []byte
}

func createTests(id string) []test {
	id += "/"
	startBA, _ := json.Marshal(&jsonCharacter{
		"tester", "25", "25", "25", "test",
	})
	return []test{
		{
			"happy start",
			startEP + id,
			http.StatusOK,
			startBA,
		},
	}
}
