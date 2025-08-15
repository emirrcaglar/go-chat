package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/emirrcaglar/go-chat/config"
	routes "github.com/emirrcaglar/go-chat/routes"
	"github.com/emirrcaglar/go-chat/server"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	mux := http.NewServeMux()
	roomStore := routes.NewRoomStore()
	s := server.NewServer(roomStore)

	routeHandler := routes.NewHandler(roomStore)
	routeHandler.RegisterRoutes(mux, s)

	portStr := strconv.Itoa(cfg.Server.Port)
	fmt.Printf("Server starting on port %s...\n", portStr)
	if err := http.ListenAndServe(":
