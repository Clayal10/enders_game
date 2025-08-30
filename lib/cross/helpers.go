package cross

import (
	"net"
)

func GetFreePort() uint16 {
	listener, _ := net.Listen("tcp", ":0")
	defer listener.Close()

	addr := listener.Addr().(*net.TCPAddr)
	return uint16(addr.Port)
}
