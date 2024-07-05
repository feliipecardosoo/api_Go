package server

import (
	"bd/src/requisicoes"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Server subo o servido
func Server() {
	router := mux.NewRouter()

	router.HandleFunc("/usuarios", requisicoes.CriarUsuario).Methods(http.MethodPost)

	fmt.Println("Escutando na porta 5k")
	log.Fatal(http.ListenAndServe(":5000", router))
}
