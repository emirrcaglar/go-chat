package types

import (
	"sync"
	"time"
)

type RoomStore struct {
	Rooms map[int]*Room
}

type Room struct {
	RoomIndex      int
	MessageHistory []Message
	roomMutex      sync.Mutex
}

func (r *Room) AddMessage(username, content string) {
	r.roomMutex.Lock()
	defer r.roomMutex.Unlock()

	r.MessageHistory = append(r.MessageHistory, Message{
		Username:  username,
		Content:   content,
		Timestamp: time.Now(),
	})
}
