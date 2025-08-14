package routes

import (
	"log"
	"net/http"
)

func (h *Handler) logoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		session, _ := store.Get(r, "sess")
		session.Values["username"] = ""

		if err := session.Save(r, w); err != nil {
			log.Printf("Error saving session: %v", err)
			http.Error(w, "Session error", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
