package server

type Config struct {
	Port uint16 `json:"ServerPort"`
}

// New will create a new server instance that starts all necessary processes
// for the server. The function will return a list of functions that should be called
// for the termination of the server
func New(cfg *Config) ([]func(), error) {
	game := newGame()

	rec, err := newReceiver(cfg, game)
	if err != nil {
		return nil, err
	}

	rec.start()

	return []func(){
		rec.stop,
	}, nil
}
