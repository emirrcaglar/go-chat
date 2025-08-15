package routes

import (
	"log"
	"net/http"
)

func (h *Handler) logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "sess")
	// Debug: Print what's in the session before clearing
	log.Printf("Before logout - session values: %+v", session.Values)

	delete(session.Values, "username")

	// Debug: Print what's in the session after clearing
	log.Printf("After logout - session values: %+v", session.Values)
	if err := session.Save(r, w); err != nil {
		log.Printf("Error saving session: %v", err)
		http.Error(w, "Session error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
