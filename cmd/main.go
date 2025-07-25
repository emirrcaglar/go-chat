package main

import (
	"fmt"
	"net/http"

	"github.com/emirrcaglar/go-chat/route"
	"github.com/emirrcaglar/go-chat/server"
	"github.com/gorilla/mux"
	"golang.org/x/net/websocket"
)

func main() {
	server := server.NewServer()

	router := mux.NewRouter()
	subrouter := router.PathPrefix("/room").Subrouter()

	http.Handle("/ws", websocket.Handler(server.HandleWS))

	routeHandler := route.NewHandler()
	routeHandler.RegisterRoutes(router, subrouter)

	fmt.Println("Server starting on port 3000...")
	if err := http.ListenAndServe(":3000", router); err != nil {
		fmt.Printf("Error starting server:%v", err)
	}
}
