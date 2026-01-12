package main

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/gorilla/websocket"
)

// Una conexion de websocket es una conexion de tipo TCP; la conexion inicia
// haciendo una peticion HTTP Get con la diferencia de mandar una cabecera o header
// denominada `Upgrade: websocket`. Una vez que el servidor acepta esta solicitud, se inicia
// una conexcion o canal persistente.

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	http.HandleFunc("/ws", handleWS)
	fmt.Println("Servidor WebSocket en http://localhost:8080/ws")
	http.ListenAndServe(":8080", nil)
}

func handleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("error on upgrading", err)
		return
	}

	defer conn.Close()

	fmt.Println("Client connected")

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	paused := true
	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println("client just disconnected", err)
				return
			}
			switch string(msg) {
			case "paused":
				paused = true
			case "pause":
				paused = true
				conn.WriteJSON(map[string]string{"status": "paused"})
			case "resume":
				paused = false
				conn.WriteJSON(map[string]string{"status": "resumed"})
			default:
				conn.WriteJSON(map[string]string{"error": "comando no reconocido"})
			}
		}

	}()

	for {
		select {
		case <-done:
			log.Println("Stopping sender loop")
			return
		case <-ticker.C:
			if !paused {
				var m runtime.MemStats
				runtime.ReadMemStats(&m)
				data := map[string]any{
					"time":       time.Now().Format("15:04:05"),
					"goroutines": runtime.NumGoroutine(),
					"alloc_mb":   float64(m.Alloc) / 1024 / 1024,
				}
				if err := conn.WriteJSON(data); err != nil {
					return
				}
			}
		}
	}

}
