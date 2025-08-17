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

// newRoomFormHandler – now uses layout and passes PageData for navbar highlighting
func (h *Handler) newRoomFormHandler(w http.ResponseWriter, r *http.Request) {
	data := types.PageData{
		PageTitle:   "Create New Room",
		CurrentPage: "new-room", // highlight navbar link
	}

	err := h.roomTemplate.ExecuteTemplate(w, "room.html", struct {
		types.PageData
	}{data})
	if err != nil {
		http.Error(w, "Failed to render page: "+err.Error(), http.StatusInternalServerError)
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

	var username string
	if val, exists := session.Values["username"]; exists {
		if uname, ok := val.(string); ok && uname != "" {
			username = uname
		}
	}
	if username == "" {
		username = "unnamed"
	}

	// CHANGED: Create a combined data struct for layout
	// data := struct {
	// 	types.PageData
	// 	*types.Room
	// 	Username string
	// }{
	// 	PageData: types.PageData{
	// 		PageTitle:   "Room " + strconv.Itoa(room.RoomIndex),
	// 		CurrentPage: "", // no navbar highlight here
	// 	},
	// 	Room:     room,
	// 	Username: username,
	// }

	data := types.RoomPageData{
		PageData: types.PageData{
			PageTitle:   "Room " + strconv.Itoa(room.RoomIndex),
			CurrentPage: "",
		},
		Username:       username,
		Room:           room,
		MessageHistory: room.MessageHistory,
	}

	err = h.roomTemplate.ExecuteTemplate(w, "room.html", data)
	if err != nil {
		http.Error(w, "Failed to render page: "+err.Error(), http.StatusInternalServerError)
	}
}
