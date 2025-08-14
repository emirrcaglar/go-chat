package routes

import (
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/emirrcaglar/go-chat/types"
)

var (
	roomCounter int
	roomMutex   sync.Mutex
)

func NewRoomStore() *types.RoomStore {
	return &types.RoomStore{
		Rooms: make(map[int]*types.Room),
	}
}

func NewRoom() *types.Room {
	roomMutex.Lock()
	defer roomMutex.Unlock()

	roomCounter++
	return &types.Room{
		RoomIndex: roomCounter,
	}
}

func (h *Handler) newRoomFormHandler(w http.ResponseWriter, r *http.Request) {
	data := types.PageData{
		PageTitle:   "New Room",
		CurrentPage: "new-room",
	}
	err := h.templates.ExecuteTemplate(w, "room", data)
	if err != nil {
		http.Error(w, "Failed to render page: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) createRoomHandler(w http.ResponseWriter, r *http.Request) {
	room := NewRoom()
	h.roomStore.Rooms[room.RoomIndex] = room

	http.Redirect(w, r, "/rooms/"+strconv.Itoa(room.RoomIndex), http.StatusSeeOther)
}

func (h *Handler) viewRoomHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid room ID", http.StatusBadRequest)
		return
	}

	room, exists := h.roomStore.Rooms[id]
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
		*types.Room
		Username string
	}{
		Room:     room,
		Username: username,
	}

	err = h.templates.ExecuteTemplate(w, "room", data)
	if err != nil {
		http.Error(w, "Failed to render page: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
