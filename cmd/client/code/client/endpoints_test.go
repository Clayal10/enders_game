package client

import (
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

	id := 123
	client := newClient()

}
