package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

const defaultPort = 8080
const staticDir = "./ui" // exe must be in root of repo

func main() {
	if err := serve(); err != nil {
		log.Fatal(err)
	}
}

func serve() error {
	http.HandleFunc("/", mainPageHandler)
	fs := http.FileServer(http.Dir(staticDir))
	http.Handle("/ui/", http.StripPrefix("/ui/", fs))

	return http.ListenAndServe(fmt.Sprintf(":%v", defaultPort), nil)
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
