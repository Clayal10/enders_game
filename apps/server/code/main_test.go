package main

import (
	"testing"

	"github.com/Clayal10/enders_game/lib/assert"
)

func TestConfig(t *testing.T) {
	a := assert.New(t)
	t.Run("TestNoConfigFile", func(_ *testing.T) {
		cfg, err := getServerConfig("Nope")
		a.NoError(err)
		a.True(cfg.Port == defaultPort)
	})
}
