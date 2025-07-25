package routes

import (
	"net/http"
	"sync"

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
	h.templates.ExecuteTemplate(w, "user.html", nil)
}

func (h *Handler) createUserHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "sess")

	if r.Method == "POST" {
		userName := r.FormValue("uname")

		userMutex.Lock()
		defer userMutex.Unlock()
		session.Values["username"] = userName
		session.Save(r, w)

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

}
