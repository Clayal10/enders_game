package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/Clayal10/enders_game/lib/lurk"
)

type ServerConfig struct {
	Hostname string
	Port     string
}

// Each field correlates to a section of the UI that should be updated in the JS.
// This struct will be json'd and each field should be added to the innerHTML of the HTML element,
// and should not overwrite the data already in it.
type ClientUpdate struct {
	Info    string `json:"info"`
	Rooms   string `json:"rooms"`
	Players string `json:"players"`
}

// Client contains all data needed to run a client instance.
type Client struct {
	// Message queue?

	Game *lurk.Game

	mu   sync.Mutex
	conn net.Conn
}

const defaultPort = 5068
const staticDir = "./ui" // exe must be in root of repo

func Start() error {
	registerEndpoints()
	return serve()
}

const setupEP = "/lurk-client/setup/"
const startEP = "/lurk-client/start/" // For now, this will just auto fill with default stuff.
const updateEP = "/lurk-client/update/"

func registerEndpoints() {
	http.HandleFunc(setupEP, handleSetup)
	http.HandleFunc(startEP, handleStart)
}

func handleSetup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		return
	}
	ba, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cfg := &ServerConfig{}
	if err := json.Unmarshal(ba, cfg); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, clientUpdate, err := setup(cfg)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(clientUpdate)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err = w.Write(jsonData); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func handleStart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		return
	}
}

func serve() error {
	http.HandleFunc("/", mainPageHandler)
	fs := http.FileServer(http.Dir(staticDir))
	http.Handle("/ui/", http.StripPrefix("/ui/", fs))

	return http.ListenAndServe(fmt.Sprintf(":%v", defaultPort), nil)
}

func mainPageHandler(w http.ResponseWriter, req *http.Request) {
	template, err := template.ParseFiles(fmt.Sprintf("%v/html/landing.html", staticDir))
	if err != nil {
		log.Printf("%v: could not parse HTML file", err)
		return
	}

	if err = template.Execute(w, nil); err != nil {
		log.Printf("%v: could not execute HTML file", err)
		return
	}
}

func setup(cfg *ServerConfig) (*Client, *ClientUpdate, error) {
	return connectToServer(cfg)
}

func connectToServer(cfg *ServerConfig) (*Client, *ClientUpdate, error) {
	conn, err := net.Dial("tcp", cfg.Hostname+":"+cfg.Port)
	if err != nil {
		return nil, nil, err
	}

	lurkMessages, err := func(conn net.Conn) (messages []lurk.LurkMessage, _ error) {
		for {
			conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
			ba, _, err := lurk.ReadSingleMessage(conn)
			if err != nil {
				if errors.Is(err, os.ErrDeadlineExceeded) {
					return messages, nil
				}
				return nil, err
			}

			lmsg, err := lurk.Unmarshal(ba)
			if err != nil {
				return nil, err
			}

			messages = append(messages, lmsg)
		}
	}(conn)
	if err != nil {
		return nil, nil, err
	}

	c := &Client{
		conn: conn,
	}

	cu := c.getClientUpdate(lurkMessages)

	return c, cu, nil
}

func (client *Client) getClientUpdate(lurkMessages []lurk.LurkMessage) *ClientUpdate {
	cu := &ClientUpdate{}
	for _, msg := range lurkMessages {
		switch msg.GetType() {
		case lurk.TypeGame:
			game := msg.(*lurk.Game)
			client.Game = game
			cu.Info += fmt.Sprintf("Stat Limit: %v\nInitial Points: %v\n%s\n", game.StatLimit, game.InitialPoints, game.GameDesc)
		case lurk.TypeVersion:
			version := msg.(*lurk.Version)
			extensionBytes := 0
			for _, v := range version.Extensions {
				extensionBytes += len(v)
			}
			cu.Info += fmt.Sprintf("LURK Version %v.%v | %v Bytes of Extensions\n", version.Major, version.Minor, extensionBytes)
		case lurk.TypeCharacter:
			character := msg.(*lurk.Character)
			cu.Players += fmt.Sprintf("%s | Attack: %v Defense: %v Health: %v Monster?: %v\n", character.Name, character.Attack, character.Defense, character.Health, character.Flags[lurk.Monster])
		}
	}
	return cu
}
