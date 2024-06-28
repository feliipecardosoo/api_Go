package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Server subo o servidor
func Server() {
	router := mux.NewRouter()

	fmt.Println("Escutando na porta 5k")
	log.Fatal(http.ListenAndServe(":5000", router))
}
