package main

import (
	"fmt"
	l "log"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/websocket"
)

const SERVER_PORT = "52694"

var log = l.New(os.Stderr, "", 0)

var CLIENTS sync.Map // sync.Map<string, *Client>

// Note: Gorilla's websockets support only one concurrent reader and one
// concurrent writer, and require external synchronization.
type Client struct {
	*websocket.Conn
	sync.Mutex
}

func main() {
	addr := fmt.Sprintf(":%v", SERVER_PORT)
	http.ListenAndServe(addr, http.HandlerFunc(handleRequest))
}

func handleRequest(rew http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/broadcast":
		notifyClients()

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
		log.Printf("failed to init connection at %v: %v", req.RemoteAddr, err)
		return
	}

	key := req.RemoteAddr
	CLIENTS.Store(key, &Client{Conn: conn})
	defer CLIENTS.Delete(key)

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			return
		}
	}
}

func skipOriginCheck(*http.Request) bool { return true }

func notifyClients() {
	CLIENTS.Range(func(_, val interface{}) bool {
		go notifyClient(val.(*Client))
		return true
	})
}

func notifyClient(client *Client) {
	client.Lock()
	defer client.Unlock()

	err := client.WriteMessage(websocket.TextMessage, nil)
	if err != nil {
		log.Printf("failed to notify socket: %+v", err)
	}
}
