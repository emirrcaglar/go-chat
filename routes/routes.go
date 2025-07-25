package route

import (
	"html/template"
	"net/http"
)

type Handler struct {
	roomStore *RoomStore
	templates *template.Template
}

func NewHandler() *Handler {
	tmpl := template.Must(template.ParseGlob("templates/*.html"))
	return &Handler{
		roomStore: NewRoomStore(),
		templates: tmpl,
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /", h.index)
	mux.HandleFunc("GET /rooms/{id}", h.viewRoomHandler)
	mux.HandleFunc("POST /create-room", h.createRoomHandler)
	mux.HandleFunc("GET /rooms/new", h.newRoomFormHandler)
}
