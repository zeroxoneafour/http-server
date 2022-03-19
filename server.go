// server.go - A server implementation
// tbh I stole most of this from a GitHub gist

package http_server

import (
	"fmt"
	"log"
	"net"
)

type TCPServer struct {
	host string
	port string
}

func NewTCPServer(host, port string) *TCPServer {
	server := new(TCPServer)
	server.host = host
	server.port = port
	return server
}

func (server *TCPServer) Run(handler func(net.Conn)) {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", server.host, server.port))
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handler(conn) // the handler is defined in main http-server.go file
	}
}
