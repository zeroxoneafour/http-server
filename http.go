// http.go - The HyperText Transit Protocol, but bad because I wrote it

package http_server

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Method uint8
type Status uint

const (
	GET Method = iota // various methods
	POST
)

type HTTPMessage struct { // common between requests and responses
	version string
	headers map[string]string
	content string
}

type HTTPRequest struct {
	HTTPMessage
	method Method
	uri    string // the location of the resource
}

type HTTPResponse struct {
	HTTPMessage
	status     Status
	statusText string // text like "Not Found"
}

type HTTPResponseConfig struct { // default res settings
	statuses map[Status]string
	headers  map[string]string
}

func NewHTTPRequest() *HTTPRequest {
	ret := new(HTTPRequest)
	ret.headers = make(map[string]string)
	return ret
}

func (r *HTTPRequest) ReadRequest(req io.Reader) error {
	reader := bufio.NewReader(req)
	firstLine := true // first line shows what type request is
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return nil // return at eof (probably never going to be reached in networking)
		}
		if firstLine {
			components := strings.Fields(line) // ex. GET / HTTP/1.1
			if len(components) != 3 {
				return errors.New("Failed reading HTTP Method")
			}
			switch components[0] { // ex. GET
			case "GET":
				r.method = GET
			case "POST":
				r.method = POST
			default:
				return errors.New("Failed reading HTTP Method")
			}
			r.uri = components[1]     // ex. /
			r.version = components[2] // ex. HTTP/1.1
			firstLine = false
		} else { // basic headers and stuff
			header, value, found := strings.Cut(line, ":") // cut at ':'. Function below trims off spaces and stuff
			if found {
				r.headers[strings.Title(strings.Trim(header, " \n\r"))] = strings.Trim(value, " \n\r") // Splits ex. Host: localhost:8000 into { "Host": "localhost:8000" }
			} else {
				if value, ok := r.headers["Content-Length"]; ok { // if no content then just don't do the parsing
					contentLength, _ := strconv.Atoi(value)
					buffer := make([]byte, contentLength)
					reader.Read(buffer) // read the content
					r.content = string(buffer)
				}
				return nil
			}
		}
	}
}

func NewHTTPResponse(defaults *HTTPResponseConfig) *HTTPResponse {
	ret := new(HTTPResponse)
	ret.version = "HTTP/1.1" // can't really change this right now, not compliant with ex. HTTP/2 really
	ret.headers = make(map[string]string)
	if defaults != nil {
		for header, value := range defaults.headers {
			ret.headers[header] = value
		}
	}
	return ret
}

func (r *HTTPResponse) SetStatus(status Status, defaults *HTTPResponseConfig) { // to be invoked outside of user control
	r.status = status
	if defaults != nil {
		r.statusText = defaults.statuses[status]
	}
}

func (r *HTTPResponse) String() string { // converts a *HTTPResponse to a string for sending to client
	ret := r.version + " " + fmt.Sprint(r.status) + " " + string(r.statusText) + "\n"
	for header, value := range r.headers {
		ret += header + ": " + value + "\n"
	}
	if len(r.content) > 0 {
		ret += "\n" // extra \n because http
		ret += r.content
	}
	return ret
}

func NewHTTPResponseConfig() *HTTPResponseConfig { // this config is for default headers and stuff so they aren't required to be set manually each time
	ret := new(HTTPResponseConfig)
	ret.statuses = make(map[Status]string)
	ret.headers = make(map[string]string)
	// default statuses
	ret.statuses[200] = "OK"
	ret.statuses[201] = "Created"
	ret.statuses[403] = "Forbidden"
	ret.statuses[404] = "Not Found"
	ret.statuses[418] = "I'm a teapot" // teapot
	// go server
	ret.headers["Server"] = "Go http-server"
	return ret
}
