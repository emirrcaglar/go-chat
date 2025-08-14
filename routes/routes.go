package routes

import (
	"html/template"
	"net/http"

	"github.com/emirrcaglar/go-chat/server"
	"github.com/emirrcaglar/go-chat/types"
	"golang.org/x/net/websocket"
)

type Handler struct {
	roomStore *types.RoomStore
	templates *template.Template
}

func NewHandler() *Handler {
	tmpl := template.Must(template.ParseGlob("templates/*.html"))
	return &Handler{
		roomStore: NewRoomStore(),
		templates: tmpl,
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux, s *server.Server) {
	mux.Handle("GET /ws/room/", websocket.Handler(s.HandleWS))

	mux.HandleFunc("GET /", h.index)
	mux.HandleFunc("GET /rooms/{id}", h.viewRoomHandler)
	mux.HandleFunc("GET /rooms/new", h.newRoomFormHandler)
	mux.HandleFunc("GET /username", h.usernameHandler)

	mux.HandleFunc("POST /create-room", h.createRoomHandler)
	mux.HandleFunc("POST /create-user", h.createUserHandler)
	mux.HandleFunc("POST /logout", h.logoutHandler)
}
