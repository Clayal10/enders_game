package cross

import (
	"log"
	"net"
)

func GetFreePort() uint16 {
	listener, _ := net.Listen("tcp", ":0")
	addr := listener.Addr().(*net.TCPAddr)
	_ = listener.Close()
	return uint16(addr.Port)
}

func LogOnErr(f func() error) {
	if err := f(); err != nil {
		log.Printf("%v: error in deferred function", err.Error())
	}
}
