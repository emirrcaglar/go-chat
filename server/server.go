package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strconv"
	"sync"

	"github.com/emirrcaglar/go-chat/types"
	"golang.org/x/net/websocket"
)

type Server struct {
	conns     map[*websocket.Conn]bool
	roomStore *types.RoomStore
	mutex     sync.RWMutex
}

func NewServer(roomStore *types.RoomStore) *Server {
	return &Server{
		conns:     make(map[*websocket.Conn]bool),
		roomStore: roomStore,
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
	// slice of bytes with size of 1024
	buf := make([]byte, 1024)

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

		id, err := strconv.Atoi(msg.RoomID)
		if err != nil {
			log.Printf("Invalid room ID: %s", msg.RoomID)
			return
		}

		if room, exists := s.roomStore.Rooms[id]; exists {
			room.AddMessage(msg.Username, msg.Content)
			log.Printf("DEBUG - Message history: %v", room.MessageHistory)
		} else {
			log.Printf("DEBUG - Message history: %v", room.MessageHistory)
			log.Printf("Room %d not found\n", id)
		}

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
