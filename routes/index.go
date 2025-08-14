package routes

import (
	"log"
	"net/http"

	"github.com/emirrcaglar/go-chat/types"
)

func (h *Handler) index(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "sess")
	if err != nil {
		log.Printf("error getting session: %v\n", err)
	}
	username := "unnamed" // default value

	if val, exists := session.Values["username"]; exists {
		if uname, ok := val.(string); ok {
			username = uname
		}
	}

	data := types.IndexPageData{
		PageData: types.PageData{
			PageTitle:   "Chat",
			CurrentPage: "index",
		},
		Rooms:    h.roomStore.Rooms,
		Username: username,
	}

	err = h.templates.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		http.Error(w, "Failed to render page: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
