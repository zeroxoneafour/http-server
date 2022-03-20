package http_server

import (
	"fmt"
	"os"
	"testing"
)

func getHandler(client *HTTPClient) Status {
	filename := "." + client.Req.GetPath()
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		return 404
	}
	fileinfo, _ := os.Stat(filename)
	content, err := os.ReadFile(filename)
	if err != nil {
		return 403
	}
	client.Res.Content = string(content)
	client.Res.Headers["Content-Length"] = fmt.Sprint(fileinfo.Size())
	return 200
}

func TestHTTPServer(t *testing.T) {
	server := New("localhost", "8000")
	server.SetHandler(GET, getHandler)
	server.Defaults.Headers["Server"] = "Go test server"
	server.Run()
}
