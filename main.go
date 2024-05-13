package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/EduardoTLopes/fernanda/database"
	"github.com/EduardoTLopes/fernanda/routes"
)

func main() {
	database.DatabaseConnection()
	handler := routes.Routes()
	port := ":8082"
	fmt.Printf("Server running on port %s", port)
	log.Fatal(http.ListenAndServe(port, handler))
}
