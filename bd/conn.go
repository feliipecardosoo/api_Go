package bd

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// ConexaoBD faz a minha conexao com o banco de dados
func ConexaoBD() (*sql.DB, error) {
	stringconn := "root:3103@/devbook?charset=utf8&parseTime=True&loc=Local"

	db, erro := sql.Open("mysql", stringconn)
	if erro != nil {
		log.Fatal(erro)
		return nil, erro
	}

	if erro = db.Ping(); erro != nil {
		log.Fatal(erro)
		return nil, erro
	}

	return db, nil
}
