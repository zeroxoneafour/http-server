// http-server.go - Bringing it all together

package http_server

import (
	"log"
	"net"
)

type HTTPServer struct { // main server, contains resDefaults config, a tcp server, and handlers for various methods
	server      *TCPServer
	resDefaults *HTTPResponseConfig
	handlers    []func(*HTTPClient) Status // returns status so it can be set manually by server
}

type HTTPClient struct { // a client, with raw IO and request/response structs
	conn net.Conn
	req  *HTTPRequest
	res  *HTTPResponse
}

func New(host, port string) *HTTPServer { // creates new HTTPServer obviously
	ret := new(HTTPServer)
	ret.server = new(TCPServer)
	ret.server.host = host
	ret.server.port = port
	ret.handlers = make([]func(*HTTPClient) Status, 9) // GET HEAD POST PUT DELETE CONNECT OPTIONS TRACE PATCH
	ret.resDefaults = NewHTTPResponseConfig()
	return ret
}

func (s *HTTPServer) SetHandler(method Method, handler func(*HTTPClient) Status) { // sets handler for operation
	s.handlers[method] = handler
}

func (s *HTTPServer) handleRequest(conn net.Conn) {
	client := new(HTTPClient) // new HTTPClient
	client.conn = conn        // not really neccessary to the struct and probably a security risk, oh well
	client.req = NewHTTPRequest()
	err := client.req.ReadRequest(client.conn) // read request from connection raw data
	if err != nil {
		log.Fatal(err) // idk
	}
	client.res = NewHTTPResponse(s.resDefaults) // initialize with defaults
	if s.handlers[client.req.method] != nil {
		status := s.handlers[client.req.method](client) // invoke handler for method, set by s.SetHandler
	} else {
		status := 501 // not implemented
	}
	client.res.SetStatus(status, s.resDefaults)    // set statusString with defaults
	client.conn.Write([]byte(client.res.String())) // write response to client
	client.conn.Close()                            // pretty important
}

func (s *HTTPServer) Run() {
	s.server.Run(s.handleRequest) // runs the TCPServer with the requesthandler defined in the server, basically just wraps the code I stole with HTTP parsing
}
