package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/Clayal10/enders_game/apps/server/code/server"
)

const (
	defaultPort     = 34567
	defaultFilePath = "./Config.json"
)

var (
	defaultConfig = fmt.Sprintf(`{
		"ServerPort": %v,
	}`, defaultPort)
)

func main() {
	cfg, err := getServerConfig(defaultFilePath)
	fatalOnErr(err)

	cancelFunctions := server.New(cfg)

	for _, cf := range cancelFunctions {
		cf()
	}
}

func getServerConfig(path string) (*server.ServerConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return nil, err
		}
		data = []byte(defaultConfig)
	}

	cfg := &server.ServerConfig{}
	if err = json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func fatalOnErr(err error) {
	if err != nil {
		log.Fatalf("%v: Could not start server!", err.Error())
	}
}
