package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/emirrcaglar/go-chat/types"
	"golang.org/x/net/websocket"
)

type Server struct {
	conns     map[*websocket.Conn]bool
	rooms     map[int]map[*websocket.Conn]bool
	roomStore *types.RoomStore
	mutex     sync.RWMutex
}

func NewServer(roomStore *types.RoomStore) *Server {
	return &Server{
		conns:     make(map[*websocket.Conn]bool),
		rooms:     make(map[int]map[*websocket.Conn]bool),
		roomStore: roomStore,
	}
}

func (s *Server) HandleWS(ws *websocket.Conn) {
	fmt.Printf("new connection: %v\n", ws.RemoteAddr())

	roomIDStr := strings.TrimPrefix(ws.Request().URL.Path, "/ws/room/")
	roomID, err := strconv.Atoi(roomIDStr)
	if err != nil {
		log.Printf("Invalid room ID: %s", roomIDStr)
		return
	}

	// create connection for the room
	s.mutex.Lock()
	if _, ok := s.rooms[roomID]; !ok {
		s.rooms[roomID] = make(map[*websocket.Conn]bool)
	}
	s.rooms[roomID][ws] = true
	s.mutex.Unlock()

	defer func() {
		s.mutex.Lock()
		delete(s.rooms[roomID], ws)
		if len(s.rooms[roomID]) == 0 {
			delete(s.rooms, roomID)
		}
		s.mutex.Unlock()
		log.Printf("connection closed: %v\n", ws.RemoteAddr())
	}()

	s.readLoop(ws)
}

func (s *Server) readLoop(ws *websocket.Conn) {
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
			// Set the timestamp on the server side
			msg.Timestamp = time.Now()
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
		s.broadcast(id, msgBytes)
	}
}

func (s *Server) broadcast(roomID int, msg []byte) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if conns, ok := s.rooms[roomID]; ok {
		for conn := range conns {
			if _, err := conn.Write(msg); err != nil {
				log.Printf("Broadcast error to %v: %v", conn.RemoteAddr(), err)
			}
		}
	}
}
