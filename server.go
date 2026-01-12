// server.go
package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("📨 Recibí:", r.Method, r.URL.Path)
	w.Header().Set("Connection", "close") // forzamos cierre TCP
	fmt.Fprintf(w, "Hola desde servidor\n")
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Servidor en http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
