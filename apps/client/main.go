package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"

	"github.com/Clayal10/enders_game/apps/client/client"
)

const defaultPort = 5068
const staticDir = "./ui" // exe must be in root of repo

func main() {
	registerEndpoints()
	if err := serve(); err != nil {
		log.Fatal(err)
	}
}

const setupEP = "/lurk-client/setup/"

func registerEndpoints() {
	http.HandleFunc(setupEP, handleSetup)
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

	cfg := &client.ServerConfig{}
	if err := json.Unmarshal(ba, cfg); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	client, err := client.Setup(cfg)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(client.Game)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err = w.Write(jsonData); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	go client.Start()
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
