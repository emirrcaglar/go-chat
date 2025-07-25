package route

import (
	"net/http"
)

type IndexPageData struct {
	PageTitle string
}

func (h *Handler) index(w http.ResponseWriter, r *http.Request) {
	data := IndexPageData{
		PageTitle: "Chat",
	}

	h.templates.ExecuteTemplate(w, "index.html", data)
}
