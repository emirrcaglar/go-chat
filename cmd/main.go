package main

import (
	"fmt"
	"net/http"

	"github.com/emirrcaglar/go-chat/cmd/server"
	"golang.org/x/net/websocket"
)

func main() {
	server := server.NewServer()

	http.Handle("/ws", websocket.Handler(server.handleWS))
	fmt.Println("Server starting on port 3000...")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		fmt.Printf("Error starting server:%v", err)
	}
}
