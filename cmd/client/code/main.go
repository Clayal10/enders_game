package main

import (
	"encoding/json"
	"fmt"
	"html/template"
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
	serve()
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

func serve() error {
	http.HandleFunc("/", mainPageHandler)
	fs := http.FileServer(http.Dir(staticDir))
	http.Handle("/ui/", http.StripPrefix("/ui/", fs))

	return http.ListenAndServe(fmt.Sprintf("0.0.0.0:%v", defaultPort), nil)
}

func mainPageHandler(w http.ResponseWriter, req *http.Request) {
	template, err := template.ParseFiles(fmt.Sprintf("%v/html/home.html", staticDir))
	if err != nil {
		log.Printf("%v: could not parse HTML file", err)
		return
	}

	if err = template.Execute(w, nil); err != nil {
		log.Printf("%v: could not execute HTML file", err)
		return
	}
}
