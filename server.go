// server.go
package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(519)
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Servidor en http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
