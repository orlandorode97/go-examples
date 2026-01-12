// client.go
package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	resp, err := http.Get("http://localhost:8080/")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close() // 👈 este cierre es lo que genera el FIN

	io.Copy(io.Discard, resp.Body)
	fmt.Println("✅ Petición terminada")
}
