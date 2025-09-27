package client

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/Clayal10/enders_game/apps/server/code/server"
	"github.com/Clayal10/enders_game/lib/assert"
	"github.com/Clayal10/enders_game/lib/cross"
)

func TestStartingServer(t *testing.T) {
	a := assert.New(t)
	serverPort := cross.GetFreePort()

	serverConfig := &server.ServerConfig{
		Port: serverPort,
	}

	cfs, err := server.New(serverConfig)
	a.NoError(err)
	defer func() {
		for _, cf := range cfs {
			cf()
		}
	}()

	clientConfig := &Config{
		Port: fmt.Sprint(serverPort),
	}

	client, cu, err := New(clientConfig)
	a.NoError(err)
	a.True(strings.Contains(cu.Info, "LURK")) // from the version
	a.True(client.id != 0)
	a.True(client.id == cu.Id)

	client.Start()
	defer client.cf()

	time.Sleep(time.Second)

}
