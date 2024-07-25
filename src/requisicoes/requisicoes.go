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

type Usuario struct {
	Nome  string `json:"nome"`
	Email string `json:"email"`
}

func CriarUsuario(w http.ResponseWriter, r *http.Request) {
	// Resgatando body
	bodyReq, erro := ioutil.ReadAll(r.Body)
	if erro != nil {
		w.Write([]byte("Falha ao ler o corpo da requisicao"))
		return
	}

	var usuario Usuario

	// Converter os valores de byte para struct
	if erro = json.Unmarshal(bodyReq, &usuario); erro != nil {
		w.Write([]byte("Erro ao converter o usuario para struct"))
		return
	}

	// Fazendo conexao banco de dados
	db, erro := bd.ConexaoBD()
	if erro != nil {
		w.Write([]byte("Erro ao conectar com o banco de dados"))
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

func RetornarAll(w http.ResponseWriter, r *http.Request) {
	// Conectar no banco de dados
	db, err := bd.ConexaoBD()
	if err != nil {
		w.Write([]byte("Erro ao conectar no banco de dados"))
		return
	}
	defer db.Close()

	// Executar a Query
	rows, err := db.Query("SELECT * FROM usuarios")
	if err != nil {
		w.Write([]byte("Erro ao executar query no banco de dados"))
		return
	}
	defer rows.Close()

	var usuarios []usuario

	// Iterar sobre os resultados
	for rows.Next() {
		var usuario usuario
		if err := rows.Scan(&usuario.ID, &usuario.Nome, &usuario.Email); err != nil {
			w.Write([]byte("Erro ao escanear os resultados"))
			return
		}
		usuarios = append(usuarios, usuario)
	}

	if err = rows.Err(); err != nil {
		w.Write([]byte("Erro ao iterar os resultados"))
		return
	}

	// Retornar a resposta em JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(usuarios); err != nil {
		w.Write([]byte("Erro ao converter os resultados para JSON"))
		return
	}
}

func EditarUsuario(w http.ResponseWriter, r *http.Request) {
	// Capturar informacao do usuario do header
	vars := mux.Vars(r)
	idStr := vars["id"]
	// Convertendo para inteiro
	id, erro := strconv.Atoi(idStr)
	if erro != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	// Fazendo conexão com o banco de dados
	db, erro := bd.ConexaoBD()
	if erro != nil {
		http.Error(w, "Erro ao conectar no banco de dados", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Parseando o corpo da requisição para obter os dados a serem atualizados
	var usuario struct {
		Nome  string `json:"nome"`
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&usuario); err != nil {
		http.Error(w, "Erro ao decodificar o corpo da requisição", http.StatusBadRequest)
		return
	}

	// Fazer a query de UPDATE
	statement, erro := db.Prepare(`
		UPDATE usuarios
		SET nome = ?, email = ?
		WHERE id = ?
	`)
	if erro != nil {
		http.Error(w, "Erro ao criar statement", http.StatusInternalServerError)
		return
	}
	defer statement.Close()

	// Executar a query de UPDATE
	resultado, erro := statement.Exec(usuario.Nome, usuario.Email, id)
	if erro != nil {
		http.Error(w, "Erro ao executar a query de update", http.StatusInternalServerError)
		return
	}

	// Verificar quantas linhas foram afetadas
	numLinhasAfetadas, erro := resultado.RowsAffected()
	if erro != nil {
		http.Error(w, "Erro ao obter número de linhas afetadas", http.StatusInternalServerError)
		return
	}

	// Verificar se algum registro foi atualizado
	if numLinhasAfetadas == 0 {
		http.Error(w, "Nenhum usuário encontrado para atualizar", http.StatusNotFound)
		return
	}

	// Retornar resposta de sucesso
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Usuário atualizado com sucesso"))
}
