package routes

import (
	"log"
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
		return // Don't forget this return!
	}

	room, exists := h.roomStore.rooms[id]
	if !exists {
		http.NotFound(w, r)
		return
	}

	// Get username from session
	session, err := store.Get(r, "sess")
	if err != nil {
		log.Printf("error getting session: %v\n", err)
	}

	username := "unnamed"
	if val, exists := session.Values["username"]; exists {
		if uname, ok := val.(string); ok {
			username = uname
		}
	}

	if username == "" {
		username = "unnamed"
	}

	data := struct {
		*Room
		Username string
	}{
		Room:     room,
		Username: username,
	}

	err = h.templates.ExecuteTemplate(w, "room.html", data)
	if err != nil {
		http.Error(w, "Failed to render page: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
