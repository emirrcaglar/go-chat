package main

import (
	"fmt"
	"io"
	"net/http"

	"golang.org/x/net/websocket"
)

type Server struct {
	conns map[*websocket.Conn]bool
}

func NewServer() *Server {
	return &Server{
		conns: make(map[*websocket.Conn]bool),
	}
}

func (s *Server) handleWS(ws *websocket.Conn) {
	fmt.Printf("new connection: %v", ws.RemoteAddr())

	s.conns[ws] = true
	s.readLoop(ws)
}

func (s *Server) readLoop(ws *websocket.Conn) {
	buf := make([]byte, 1024) // slice of bytes with size of 1024

	for {
		n, err := ws.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Printf("read error:%v", err)
			continue
		}
		msg := buf[:n]
		fmt.Println(string(msg))

		ws.Write([]byte("thank you for the msg"))
	}
}

func main() {
	server := NewServer()

	http.Handle("/ws", websocket.Handler(server.handleWS))
	fmt.Println("Server starting on port 3000...")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		fmt.Printf("Error starting server:%v", err)
	}
}
