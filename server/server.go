package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"github.com/emirrcaglar/go-chat/types"
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

func (s *Server) HandleWS(ws *websocket.Conn) {
	fmt.Printf("new connection: %v\n", ws.RemoteAddr())

	defer func() {
		s.mutex.Lock()
		delete(s.conns, ws)
		s.mutex.Unlock()
		log.Printf("connection closed: %v\n", ws.RemoteAddr())
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
		var msg types.Message

		if err := json.Unmarshal(buf[:n], &msg); err != nil {
			log.Printf("JSON decode error: %v", err)
			continue
		}

		msg.Timestamp = time.Now()

		msgBytes, err := json.Marshal(msg)
		if err != nil {
			log.Printf("JSON encode error: %v", err)
			continue
		}
		s.broadcast(msgBytes)
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
