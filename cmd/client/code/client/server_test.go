package client

import (
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
	"github.com/Clayal10/enders_game/pkg/lurk"
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
		defer cross.LogOnErr(fd.Close)

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
		defer cross.LogOnErr(resp.Body.Close)
		a.True(resp.StatusCode == http.StatusOK)

		a.Eventually(func() bool {
			resp, err := http.Get(fmt.Sprintf(uri, updateEP))
			a.NoError(err)
			a.True(resp.StatusCode == http.StatusOK)
			defer cross.LogOnErr(resp.Body.Close)
			return strings.Contains(client.State.Players, "Colonel Graph")
		}, time.Second*2, time.Millisecond*100)

		bot := startClientConnection(a, serverConfig, &lurk.Character{
			Name:       "bot",
			Attack:     10,
			Defense:    20,
			Regen:      30,
			PlayerDesc: "A bot who is going to send a message",
		})

		_, err = bot.Write(lurk.Marshal(&lurk.Message{
			Recipient: client.character.Name,
			Sender:    "bot",
			Text:      "HELLO",
		}))
		a.NoError(err)

		a.Eventually(func() bool {
			resp, err := http.Get(fmt.Sprintf(uri, updateEP))
			a.NoError(err)
			a.True(resp.StatusCode == http.StatusOK)
			defer cross.LogOnErr(resp.Body.Close)
			return strings.Contains(client.State.Info, "HELLO")
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

func startClientConnection(a *assert.Assert, cfg *server.Config, char *lurk.Character) net.Conn {
	conn, err := net.Dial("tcp", fmt.Sprintf(":%v", cfg.Port))
	a.NoError(err)

	buffer, n, err := lurk.ReadSingleMessage(conn)
	a.NoError(err)

	msg, err := lurk.Unmarshal(buffer[:n])
	a.NoError(err)

	a.True(msg.GetType() == lurk.TypeVersion)

	buffer, n, err = lurk.ReadSingleMessage(conn)
	a.NoError(err)

	msg, err = lurk.Unmarshal(buffer[:n])
	a.NoError(err)

	a.True(msg.GetType() == lurk.TypeGame)

	ba := lurk.Marshal(char)

	_, err = conn.Write(ba) // send character
	a.NoError(err)

	buffer, n, err = lurk.ReadSingleMessage(conn)
	a.NoError(err)

	msg, err = lurk.Unmarshal(buffer[:n])
	a.NoError(err)

	a.True(msg.GetType() == lurk.TypeCharacter)

	buffer, n, err = lurk.ReadSingleMessage(conn) // read accept
	a.NoError(err)

	msg, err = lurk.Unmarshal(buffer[:n])
	a.NoError(err)

	a.True(msg.GetType() == lurk.TypeAccept)

	ba = lurk.Marshal(&lurk.Start{
		Type: lurk.TypeStart,
	})

	_, err = conn.Write(ba)
	a.NoError(err)

	buffer, _, err = lurk.ReadSingleMessage(conn)
	a.NoError(err)

	a.True(buffer[0] == byte(lurk.TypeAccept))

	return conn
}
