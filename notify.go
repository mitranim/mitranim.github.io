package main

import (
	"fmt"
	l "log"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/websocket"
)

const (
	SERVER_PORT = "52694"
	PUBLIC_DIR  = "public"
)

// sync.Map<string, *websocket.Conn>
var CONNS sync.Map

var CLEAR_TERM = []byte("\x1bc\x1b[3J")

var log = l.New(os.Stderr, "", 0)

func main() {
	http.ListenAndServe(
		fmt.Sprintf(":%v", SERVER_PORT),
		http.HandlerFunc(handleRequest),
	)
}

func handleRequest(rew http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/broadcast":
		broadcast()

	case "/ws":
		if req.Method != http.MethodGet {
			rew.WriteHeader(http.StatusMethodNotAllowed)
			break
		}
		initConn(rew, req)

	default:
		rew.WriteHeader(http.StatusNotFound)
	}
}

func initConn(rew http.ResponseWriter, req *http.Request) {
	up := websocket.Upgrader{CheckOrigin: skipOriginCheck}
	conn, err := up.Upgrade(rew, req, nil)
	if err != nil {
		log.Printf("failed to init connection at %v: %v", req.URL, err)
		return
	}

	key := req.RemoteAddr
	CONNS.Store(key, conn)
	defer CONNS.Delete(key)

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			return
		}
	}
}

func skipOriginCheck(*http.Request) bool { return true }

func broadcast() {
	CONNS.Range(func(_, val interface{}) bool {
		conn := val.(*websocket.Conn)
		err := conn.WriteMessage(websocket.TextMessage, nil)
		if err != nil {
			log.Printf("failed to notify socket: %+v", err)
		}
		return true
	})
}
