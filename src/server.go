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

	// Criacao de Usuario
	router.HandleFunc("/usuarios", requisicoes.CriarUsuario).Methods(http.MethodPost)
	// Resgatando Usuarios
	router.HandleFunc("/usuarios", requisicoes.RetornarAll).Methods(http.MethodGet)
	// Resgatando Usuario por ID
	router.HandleFunc("/usuarios/{id}", requisicoes.RetornarUsuario).Methods(http.MethodGet)
	// Editando Usuario
	router.HandleFunc("/usuarios/edit/{id}", requisicoes.EditarUsuario).Methods(http.MethodGet)

	// Deletando Usuario

	fmt.Println("Escutando na porta 5k")
	log.Fatal(http.ListenAndServe(":5000", router))
}
