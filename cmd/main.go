package main

import (
	"fmt"
	"net/http"

	routes "github.com/emirrcaglar/go-chat/routes"
	"github.com/emirrcaglar/go-chat/server"
)

func main() {
	// server := server.NewServer()

	mux := http.NewServeMux()
	s := server.NewServer()
	// http.Handle("/ws", websocket.Handler(server.HandleWS))

	routeHandler := routes.NewHandler()
	routeHandler.RegisterRoutes(mux, s)

	fmt.Println("Server starting on port 3000...")
	if err := http.ListenAndServe(":3000", mux); err != nil {
		fmt.Printf("Error starting server:%v", err)
	}
}
