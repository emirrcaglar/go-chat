package route

import (
	"net/http"
	"text/template"
)

type RoomPageData struct {
	RoomIndex int
}

func (h *Handler) Room(w http.ResponseWriter, r *http.Request) {
	data := RoomPageData{
		RoomIndex: 1,
	}

	tmpl, _ := template.ParseFiles("templates/room.html")

	tmpl.Execute(w, data)
}
