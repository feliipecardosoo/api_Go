package main

import (
	"bd/bd"
	server "bd/src"
)

func main() {
	bd.ConexaoBD()
	server.Server()
}
