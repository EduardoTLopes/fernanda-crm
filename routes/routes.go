package routes

import (
	"net/http"

	"github.com/EduardoTLopes/fernanda/controllers"
)

func Routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	mux.HandleFunc("/", controllers.Home)
	mux.HandleFunc("/place/{ID}", controllers.ShowPlace)
	mux.HandleFunc("/places/", controllers.ShowAllPlaces)
	mux.HandleFunc("/places/new/", controllers.AddPlaces)
	mux.HandleFunc("/updateContactNameAndNotes/{ID}", controllers.UpdateContactNameAndNotes)

	return mux
}
