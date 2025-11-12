package client

import (
	"bytes"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

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

		clientPort := cross.GetFreePort()
		//nolint:errcheck // In a test.
		go http.ListenAndServe(fmt.Sprintf("0.0.0.0:%v", clientPort), nil)

		filename := filepath.Join(os.TempDir(), "test.txt")
		fd, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, os.ModePerm)
		a.NoError(err)
		defer fd.Close()

		n, err := fd.Write([]byte(`{
			"name": "Tester!",
			"attack": "nil"
		}`))
		a.NoError(err)
		a.True(n > 0)

		file, err := os.Open(filename)
		a.NoError(err)

		client, err := New(clientConfig)
		a.NoError(err)
		a.True(strings.Contains(client.State.Info, "LURK")) // from the version
		a.True(client.id != 0)
		a.True(client.id == client.State.Id)
		a.NotNil(client.Game)

		client.Start()
		defer client.cf()
		time.Sleep(50 * time.Millisecond)

		uri := fmt.Sprintf("http://localhost:%d", clientPort) + "%v" + fmt.Sprintf("%d/", client.id)
		resp, err := http.Post(fmt.Sprintf(uri, startEP), "application/json", file)
		a.NoError(err)
		defer resp.Body.Close()
		a.True(resp.StatusCode == http.StatusOK)

		resp, err = http.Post(fmt.Sprintf(uri, messageEP), "application/json", bytes.NewReader())
		a.NoError(err)
		defer resp.Body.Close()
		a.True(resp.StatusCode == http.StatusOK)
		a.Eventually(func() bool {
			resp, err := http.Get(fmt.Sprintf(uri, updateEP))
			a.NoError(err)
			a.True(resp.StatusCode == http.StatusOK)
			defer resp.Body.Close()
			return strings.Contains(client.State.Players, "Colonel Graph")
		}, time.Second*2, time.Millisecond*100)
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
