package cross_test

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
	"testing"

	"github.com/Clayal10/enders_game/lib/assert"
	"github.com/Clayal10/enders_game/lib/cross"
)

func TestGettingPort(t *testing.T) {
	a := assert.New(t)

	port := cross.GetFreePort()

	l, err := net.Listen("tcp", fmt.Sprintf("localhost:%v", port))
	a.NoError(err)
	a.True(l != nil)
}

func TestErrLogging(t *testing.T) {
	a := assert.New(t)
	var buf bytes.Buffer
	log.SetOutput(&buf)

	cross.LogOnErr(func() error {
		return errors.New("fail")
	})

	a.True(strings.Contains(buf.String(), "error in deferred function"))
}
