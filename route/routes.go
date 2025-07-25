package route

import (
	"github.com/gorilla/mux"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterRoutes(router, subrouter *mux.Router) {
	router.HandleFunc("/", h.index).Methods("GET").Name("index")
	subrouter.HandleFunc("/room", h.Room).Methods("GET")
}
