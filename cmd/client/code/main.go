package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/Clayal10/enders_game/cmd/client/code/client"
)

const setupEP = "/lurk-client/setup/"

const defaultPort = 5068
const staticDir = "../cmd/client/code/ui" // exe must be in root of repo

func main() {
	http.HandleFunc(setupEP, handleSetup)
	if err := serve(); err != nil {
		log.Fatal(err)
	}
}

func handleSetup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		return
	}
	ba, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	cfg := &client.Config{}
	if err := json.Unmarshal(ba, cfg); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	c, err := client.New(cfg)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	jsonData, err := json.Marshal(c.State)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	if _, err = w.Write(jsonData); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	c.Start()
}
