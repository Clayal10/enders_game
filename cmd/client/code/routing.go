package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"html/template"
	"io"
	"log"
	"net"
	"net/http"
)

const (
	certFile = "../../certificate/isoptera.lcsc.edu/fullchain20.pem"
	keyFile  = "../../certificate/isoptera.lcsc.edu/privkey20.pem"
)

func serve() error {
	http.HandleFunc("/", mainPageHandler)
	fs := http.FileServer(http.Dir(staticDir))

	listeningPort := fmt.Sprintf("0.0.0.0:%v", defaultPort)

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		http.Handle("/ui/", http.StripPrefix("/ui/", fs))
		log.Println("Serving over HTTP")
		return http.ListenAndServe(listeningPort, nil)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	ln, err := net.Listen("tcp", listeningPort)
	if err != nil {
		return err
	}

	log.Printf("Packet sniffing on %s", listeningPort)

	httpServer := &http.Server{
		Handler: fs,
	}
	httpsServer := &http.Server{
		Handler: fs,
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("%v: could not accept connection", err.Error())
			continue
		}
		go handlePacket(conn, tlsConfig, httpServer, httpsServer)
	}
}

func handlePacket(conn net.Conn, tlsConfig *tls.Config, httpServer, httpsServer *http.Server) {
	reader := bufio.NewReader(conn)

	b, err := reader.Peek(1)
	if err != nil {
		conn.Close()
		log.Printf("%s: could not peak in the connection", err.Error())
		return
	}
	if b[0] != 0x16 {
		log.Println(httpServer.Serve(&uniqueListener{c: &peekConn{conn, reader}}))
	}
	tlsConn := tls.Server(&peekConn{conn, reader}, tlsConfig)
	log.Println(httpsServer.Serve(&uniqueListener{c: tlsConn}))
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

type peekConn struct {
	net.Conn
	reader *bufio.Reader
}

func (p *peekConn) Read(b []byte) (int, error) {
	return p.reader.Read(b)
}

type uniqueListener struct {
	c    net.Conn
	addr net.Addr
}

func (l *uniqueListener) Accept() (net.Conn, error) {
	if l.c == nil {
		return nil, io.EOF
	}
	c := l.c
	l.c = nil
	return c, nil
}

func (l *uniqueListener) Addr() net.Addr { return l.addr }
func (*uniqueListener) Close() error     { return nil }
