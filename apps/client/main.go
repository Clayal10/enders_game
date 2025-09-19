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
		fmt.Printf("Bad Method")
		return
	}
	ba, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("Could not read")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cfg := &client.ServerConfig{}
	if err := json.Unmarshal(ba, cfg); err != nil {
		fmt.Printf("Could not unmarshal")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	client.Start(cfg)
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
