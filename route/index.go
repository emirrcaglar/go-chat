package route

import (
	"html/template"
	"net/http"
)

type IndexPageData struct {
	PageTitle string
}

func (h *Handler) index(w http.ResponseWriter, r *http.Request) {
	data := IndexPageData{
		PageTitle: "Chat",
	}

	tmpl, _ := template.ParseFiles("templates/index.html")

	tmpl.Execute(w, data)
}
