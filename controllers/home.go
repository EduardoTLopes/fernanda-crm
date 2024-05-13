package controllers

import (
	"net/http"
)

type PageData struct {
	ImageURL string
}

func Home(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		ImageURL: "/assets/fernanda.jpg",
	}

	Templates.ExecuteTemplate(w, "Home", data)
}
