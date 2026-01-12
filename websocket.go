package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var (
	Connections []map[string]*websocket.Conn
	mu          sync.Mutex
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", WebSockerHandler)
	mux.HandleFunc("/show", PrintConnections)
	mux.HandleFunc("/send", SendMessageToSpecificUser)

	log.Printf("Server stared on :8080")
	http.ListenAndServe(":8080", mux)
}

func WebSockerHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "username query param is required", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error: ", err)
		return
	}

	mu.Lock()
	Connections = append(Connections, map[string]*websocket.Conn{
		username: conn,
	})
	mu.Unlock()

	log.Printf("[%s] conneected", username)

	go handleConnection(username, conn)
}

func handleConnection(username string, conn *websocket.Conn) {
	defer func() {
		removeConnection(username)
		conn.Close()
		log.Printf("[%s] disconnected", username)
	}()

	const idleTimeout = 60 * time.Second
	conn.SetReadDeadline(time.Now().Add(idleTimeout))
	conn.SetPongHandler(func(appData string) error {
		conn.SetReadDeadline(time.Now().Add(idleTimeout))
		return nil
	})

	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("[%s] read error: %v", username, err)
			break
		}

		log.Printf("[%s] says: %s", username, msg)
		resp := fmt.Sprintf("Response to %s: %s", username, msg)
		conn.WriteMessage(msgType, []byte(resp))
	}
}

func removeConnection(username string) {
	mu.Lock()
	defer mu.Unlock()
	var updated []map[string]*websocket.Conn
	for _, m := range Connections {
		if _, ok := m[username]; !ok {
			updated = append(updated, m)
		}
	}

	Connections = updated
}

// PrintConnections displays all connected usernames
func PrintConnections(w http.ResponseWriter, _ *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	fmt.Fprintln(w, "Connected users:")
	for _, conn := range Connections {
		for username := range conn {
			fmt.Fprintln(w, "-", username)
		}
	}
}

// SendMessageToSpecificUser sends a message to a specific user
func SendMessageToSpecificUser(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	targetUser := string(body)

	mu.Lock()
	defer mu.Unlock()

	for _, item := range Connections {
		for username, conn := range item {
			if username == targetUser {
				conn.WriteMessage(websocket.TextMessage, []byte("hi "+username))
				log.Printf("Sent message to %s", username)
				w.Write([]byte("sent"))
				return
			}
		}
	}
	w.Write([]byte("user not found"))
}
