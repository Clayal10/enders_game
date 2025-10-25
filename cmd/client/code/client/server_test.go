package client

import (
	"fmt"
	"net"
	"strings"
	"testing"

	"github.com/Clayal10/enders_game/cmd/server/code/server"
	"github.com/Clayal10/enders_game/pkg/assert"
	"github.com/Clayal10/enders_game/pkg/cross"
)

func TestStartingServer(t *testing.T) {
	a := assert.New(t)
	serverPort := cross.GetFreePort()

	serverConfig := &server.Config{
		Port: serverPort,
	}

	cfs, err := server.New(serverConfig)
	a.NoError(err)
	defer func() {
		for _, cf := range cfs {
			cf()
		}
	}()

	t.Run("TestBasicSetup", func(_ *testing.T) {
		clientConfig := &Config{
			Port: fmt.Sprint(serverPort),
		}

		client, err := New(clientConfig)
		a.NoError(err)
		a.True(strings.Contains(client.State.Info, "LURK")) // from the version
		a.True(client.id != 0)
		a.True(client.id == client.State.Id)

		client.Start()
		defer client.cf()

	})
	t.Run("TestBadDial", func(_ *testing.T) {
		clientConfig := &Config{
			Port: "100000000",
		}
		_, err = New(clientConfig)
		a.Error(err)
	})
	t.Run("TestBadConnection", func(_ *testing.T) {
		badPort := cross.GetFreePort()
		l, err := net.Listen("tcp", fmt.Sprintf(":%d", badPort))
		a.NoError(err)

		go func() {
			conn, err := l.Accept()
			a.NoError(err)
			_, _ = conn.Write([]byte{0xFF, 0xFF, 0xFF, 0xFF}) // give a bad response.
		}()

		clientConfig := &Config{
			Port: fmt.Sprintf("%d", badPort),
		}
		_, err = New(clientConfig)
		a.Error(err)
	})
}
