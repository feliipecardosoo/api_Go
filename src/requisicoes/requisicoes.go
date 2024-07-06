package requisicoes

import (
	"bd/bd"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func CriarUsuario(w http.ResponseWriter, r *http.Request) {

	// Resgatando valores body
	bodyReq, erro := ioutil.ReadAll(r.Body)
	if erro != nil {
		w.Write([]byte("Falha ao ler o corpo da requisicao"))
		return
	}

	var usuario usuario

	// Converter os valores de byte para struct
	if erro = json.Unmarshal(bodyReq, &usuario); erro != nil {
		w.Write([]byte("Erro ao converter o usuario para struct"))
		return
	}

	// Fazendo conexao banco de dados
	db, erro := bd.ConexaoBD()
	if erro != nil {
		w.Write([]byte("Erro ao insert no banco de dados"))
		return
	}
	defer db.Close()

	// Preparando Statement
	statement, erro := db.Prepare("insert into usuarios (nome, email) values (?, ?)")
	if erro != nil {
		w.Write([]byte("Erro ao criar statement"))
		return
	}
	defer statement.Close()

	// Passando valores da minha query
	insercao, erro := statement.Exec(usuario.Nome, usuario.Email)
	if erro != nil {
		w.Write([]byte("Erro ao executar statement"))
		return
	}

	// Resgatando ID inserido para retorno da mensagem
	idInserido, erro := insercao.LastInsertId()
	if erro != nil {
		w.Write([]byte("Erro ao obter id inserido"))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("Usuario inserido com sucesso! Id: %d", idInserido)))
}

func RetornarUsuario(w http.ResponseWriter, r *http.Request) {

	// Obter parametros URL
	vars := mux.Vars(r)
	idStr := vars["id"]
	// Convertendo para string
	id, erro := strconv.Atoi(idStr)
	if erro != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	// Fazendo conexao banco de dados
	db, erro := bd.ConexaoBD()
	if erro != nil {
		http.Error(w, "Erro ao conectar no banco de dados", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Preparando Statement
	statement, erro := db.Prepare("SELECT * FROM usuarios WHERE id = ?")
	if erro != nil {
		http.Error(w, "Erro ao criar statement", http.StatusInternalServerError)
		return
	}
	defer statement.Close()

	var usuario usuario

	// Executando minha consulta e guardando valores no meu struct usuario
	erro = statement.QueryRow(id).Scan(&usuario.ID, &usuario.Nome, &usuario.Email)
	if erro != nil {
		if erro == sql.ErrNoRows {
			http.Error(w, "Usuário não encontrado", http.StatusNotFound)
		} else {
			http.Error(w, "Erro ao executar query", http.StatusInternalServerError)
		}
		return
	}

	// Convertendo de byte para JSON
	usuarioJSON, erro := json.Marshal(usuario)
	if erro != nil {
		http.Error(w, "Erro ao converter o usuário para JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(usuarioJSON)
}
