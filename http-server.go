// http-server.go - Bringing it all together

package http_server

import (
	"log"
	"net"
)

type HTTPServer struct { // main server, contains resDefaults config, a tcp server, and handlers for various methods
	server   *TCPServer
	Defaults *HTTPResponseConfig
	handlers []func(*HTTPClient) Status // returns status so it can be set manually by server
}

type HTTPClient struct { // a client, with raw IO and request/response structs
	conn net.Conn
	Req  *HTTPRequest
	Res  *HTTPResponse
}

func New(host, port string) *HTTPServer { // creates new HTTPServer obviously
	ret := new(HTTPServer)
	ret.server = new(TCPServer)
	ret.server.host = host
	ret.server.port = port
	ret.handlers = make([]func(*HTTPClient) Status, 9) // GET HEAD POST PUT DELETE CONNECT OPTIONS TRACE PATCH
	ret.Defaults = NewHTTPResponseConfig()
	return ret
}

func (s *HTTPServer) SetHandler(method Method, handler func(*HTTPClient) Status) { // sets handler for operation
	s.handlers[method] = handler
}

func (s *HTTPServer) handleRequest(conn net.Conn) {
	client := new(HTTPClient) // new HTTPClient
	client.conn = conn        // not really neccessary to the struct and probably a security risk, oh well
	defer client.conn.Close() // pretty important
	client.Req = NewHTTPRequest()
	err := client.Req.ReadRequest(client.conn) // read request from connection raw data
	if err != nil {
		log.Panic(err) // idk
	}
	client.Res = NewHTTPResponse(s.Defaults) // initialize with defaults
	var status Status                        // declared here so go doesn't cry
	if s.handlers[client.Req.method] != nil {
		status = s.handlers[client.Req.method](client) // invoke handler for method, set by s.SetHandler
	} else {
		status = 501 // not implemented
	}
	client.Res.SetStatus(status, s.Defaults)       // set statusString with defaults
	client.conn.Write([]byte(client.Res.String())) // write response to client
}

func (s *HTTPServer) Run() {
	s.server.Run(s.handleRequest) // runs the TCPServer with the requesthandler defined in the server, basically just wraps the code I stole with HTTP parsing
}
