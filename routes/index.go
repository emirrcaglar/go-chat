package routes

import (
	"fmt"
	"net/http"
)

type IndexPageData struct {
	PageTitle string
	Rooms     map[int]*Room
	Username  string
}

func (h *Handler) index(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "sess")
	var username = "unnamed"
	if username, exists := session.Values["username"]; exists {
		if uname, ok := username.(string); ok {
			// username exists and is a string
			username = uname
			fmt.Println("Username:", uname)
		} else {
			// username exists but isn't a string
			fmt.Println("Username is not a string")
		}
	} else {
		// username doesn't exist in session
		fmt.Println("No username in session")
	}

	data := IndexPageData{
		PageTitle: "Chat",
		Rooms:     h.roomStore.rooms,
		Username:  username,
	}

	err := h.templates.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		http.Error(w, "Failed to render page: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
