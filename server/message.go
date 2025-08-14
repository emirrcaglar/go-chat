package server

import "time"

type Message struct {
	Username  string    `json:"username"`
	Content   string    `json:"content"`
	RoomID    string    `json:"roomID"`
	Timestamp time.Time `json:"timestamp"`
}
