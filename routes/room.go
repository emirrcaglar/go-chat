package routes

import (
	"net/http"
	"strconv"
	"sync"
)

var (
	roomCounter int
	roomMutex   sync.Mutex
)

type RoomStore struct {
	rooms map[int]*Room
}

type Room struct {
	RoomIndex int
}

type Message struct {
	msgIndex int
	content  string
	sender   User
}

func NewRoomStore() *RoomStore {
	return &RoomStore{
		rooms: make(map[int]*Room),
	}
}

func NewRoom() *Room {
	roomMutex.Lock()
	defer roomMutex.Unlock()

	roomCounter++
	return &Room{
		RoomIndex: roomCounter,
	}
}

func (h *Handler) newRoomFormHandler(w http.ResponseWriter, r *http.Request) {
	h.templates.ExecuteTemplate(w, "room.html", nil)
}

func (h *Handler) createRoomHandler(w http.ResponseWriter, r *http.Request) {
	room := NewRoom()
	h.roomStore.rooms[room.RoomIndex] = room

	http.Redirect(w, r, "/rooms/"+strconv.Itoa(room.RoomIndex), http.StatusSeeOther)
}

func (h *Handler) viewRoomHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid room ID", http.StatusBadRequest)
	}

	room, exists := h.roomStore.rooms[id]
	if !exists {
		http.NotFound(w, r)
		return
	}

	h.templates.ExecuteTemplate(w, "room.html", room)
}
