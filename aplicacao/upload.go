package aplicacao

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/freehandle/breeze/crypto"
	"github.com/freehandle/breeze/protocol/actions"
)

// tamanho maximo do arquivo que pode ser publicado
const maxFileSize = 10000

type TruncatedFile struct {
	Hash  crypto.Hash
	Parts [][]byte
}

func splitBytes(bytes []byte) *TruncatedFile {
	truncated := TruncatedFile{
		Hash:  crypto.Hasher(bytes),
		Parts: make([][]byte, len(bytes)/maxFileSize+1),
	}
	for n := 0; n < len(truncated.Parts); n++ {
		if (n+1)*maxFileSize >= len(bytes) {
			truncated.Parts[n] = bytes[n*maxFileSize:]
		} else {
			truncated.Parts[n] = bytes[n*maxFileSize : (n+1)*maxFileSize]
		}

	}
	return &truncated
}

func (a *ProcuradorGeral) OperadorUpload(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(maxFileSize)
	author := a.Autor(r)
	file, _, err := r.FormFile("arquivoPraSubir")
	if err != nil {
		log.Printf("Não foi possível puxar o arquivo: %v\n", err)
		return
	}
	defer file.Close()
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		log.Printf("Erros ao ler os bytes do arquivo: %v\n", err)
	}
	var actionArray []actions.Action

	datahora := r.FormValue("dataHora")
	switch r.FormValue("acao") {
	case "PostarLivro":
		actionArray, err = FormularioLivro(r, a.estado.ArrobasPraTokens, datahora, fileBytes).ParaAcao()
	case "PostarMeme":
		actionArray, err = FormularioMeme(r, a.estado.ArrobasPraTokens, datahora, fileBytes).ParaAcao()
	}
	if err == nil && len(actionArray) > 0 {
		a.Send(actionArray, author)
	}
	http.Redirect(w, r, fmt.Sprintf("%v/", a.serverName), http.StatusSeeOther)
}
