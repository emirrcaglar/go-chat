package routes

import (
	"log"
	"net/http"
	"sync"

	"github.com/emirrcaglar/go-chat/types"
	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("what-a-secret"))

var userMutex sync.Mutex

type User struct {
	UserName string
}

func NewUser(userName string) *User {
	userMutex.Lock()
	defer userMutex.Unlock()
	return &User{
		UserName: userName,
	}
}

func (h *Handler) usernameHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "sess")
	if err != nil {
		log.Printf("error getting session: %v\n", err)
	}

	var username string
	if val, ok := session.Values["username"]; ok {
		if str, ok := val.(string); ok && str != "" {
			username = str
		}
	}

	data := types.IndexPageData{
		PageData: types.PageData{
			PageTitle:   "Set Username",
			CurrentPage: "user",
		},
		Rooms:    h.roomStore.Rooms,
		Username: username,
	}
	h.templates.ExecuteTemplate(w, "layout", data)
}

func (h *Handler) createUserHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "sess")
	if err != nil {
		log.Printf("error getting session: %v\n", err)
	}
	if r.Method == "POST" {
		userName := r.FormValue("uname")
		userMutex.Lock()
		defer userMutex.Unlock()
		session.Values["username"] = userName

		// Check for save errors
		if err := session.Save(r, w); err != nil {
			log.Printf("Error saving session: %v", err)
			http.Error(w, "Session error", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
