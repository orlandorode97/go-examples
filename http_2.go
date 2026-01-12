package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	resp, err := http.Get("https://jsonplaceholder.typicode.com/todos/1")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("Protocolo:", resp.Proto)

	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))
}
