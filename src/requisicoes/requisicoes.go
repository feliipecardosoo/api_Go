package requisicoes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func CriarUsuario(w http.ResponseWriter, r *http.Request) {

	bodyReq, erro := ioutil.ReadAll(r.Body)
	if erro != nil {
		w.Write([]byte("Falha ao ler o corpo da requisicao"))
		return
	}

	var usuario usuario

	if erro = json.Unmarshal(bodyReq, &usuario); erro != nil {
		w.Write([]byte("Erro ao converter o usuario para struct"))
		return
	}

	fmt.Println(usuario)

}
