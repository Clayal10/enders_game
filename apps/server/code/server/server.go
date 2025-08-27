package server

type ServerConfig struct {
	Port uint16 `json:"ServerPort"`
}

// New will create a new server instance that starts all necessary processes
// for the server. The function will return a list of functions that should be called
// for the termination of the server
func New(cfg *ServerConfig) []func() {

	return []func(){}
}
