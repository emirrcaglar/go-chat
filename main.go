package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"

	"golang.org/x/net/websocket"
)

type Server struct {
	conns map[*websocket.Conn]bool
	mutex sync.RWMutex
}

func NewServer() *Server {
	return &Server{
		conns: make(map[*websocket.Conn]bool),
	}
}

func (s *Server) handleWS(ws *websocket.Conn) {
	fmt.Printf("new connection: %v", ws.RemoteAddr())

	defer func() {
		s.mutex.Lock()
		delete(s.conns, ws)
		s.mutex.Unlock()
		log.Printf("connection closed: %v", ws.RemoteAddr())
	}()

	s.mutex.Lock()
	s.conns[ws] = true
	s.mutex.Unlock()

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

		if _, err := ws.Write([]byte("Echo: " + string(msg))); err != nil {
			log.Printf("write error: %v", err)
			break
		}

		s.broadcast([]byte(fmt.Sprintf("Broadcast: %s", msg)))
	}
}

func (s *Server) broadcast(msg []byte) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for conn := range s.conns {
		if _, err := conn.Write(msg); err != nil {
			log.Printf("Broadcast error to %v: %v", conn.RemoteAddr(), err)
		}
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
