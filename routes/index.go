package route

import (
	"net/http"
)

type IndexPageData struct {
	PageTitle string
	Rooms     map[int]*Room
}

func (h *Handler) index(w http.ResponseWriter, r *http.Request) {
	data := IndexPageData{
		PageTitle: "Chat",
		Rooms:     h.roomStore.rooms,
	}

	err := h.templates.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		http.Error(w, "Failed to render page: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
