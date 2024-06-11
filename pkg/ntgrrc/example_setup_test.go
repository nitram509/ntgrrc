package ntgrrc

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"testing"
)

var mockServerPort int

func TestMain(m *testing.M) {
	setupGs305EPMockServer()
	exitVal := m.Run()
	os.Exit(exitVal)
}

func setupGs305EPMockServer() {
	http.HandleFunc("/", alwaysReturn200Ok)
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}
	mockServerPort = listener.Addr().(*net.TCPAddr).Port
	log.Println(fmt.Sprintf("init GS316EP mock server on port 127.0.0.1:%d", mockServerPort))
	go http.Serve(listener, nil)
	if err != nil {
		panic(err)
	}
}

func alwaysReturn200Ok(w http.ResponseWriter, r *http.Request) {
	html := loadTestFile(string(GS316EP), "_root.html")
	_, _ = w.Write([]byte(html))
}
