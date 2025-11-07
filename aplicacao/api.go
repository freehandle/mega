package aplicacao

/*
import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/freehandle/breeze/crypto"
)

// Conversao de tipos
func FormularioParaInt(r *http.Request, campo string) int {
	if r == nil {
		return 0
	}
	valor, _ := strconv.Atoi(r.FormValue(campo))
	return valor
}

// Conversao de data
func FormularioParaDataHora(r *http.Request, campo string) time.Time {
	if r == nil {
		return time.Time{}
	}
	formato := "Mon Jan 02 2006 15:04:05 GMT-0700"
	stringData := r.FormValue(campo)
	t, katu := time.Parse(formato, stringData)
	if katu != nil {
		return t
	}
	return time.Time{}
}

// Leitura de formularios

// Leitura Formulario Causo
func FormularioCauso(r *http.Request, apelidos map[string]crypto.Token, datahora time.Time) PostarCauso {
	if r == nil {
		return PostarCauso{}
	}
	if apelidos == nil {
		return PostarCauso{}
	}
	acao := PostarCauso{
		Acao:       "PostarCauso",
		ID:         FormularioParaInt(r, "id"),
		CampoCauso: r.FormValue("campoCauso"),
		DataHora:   datahora,
	}
	return acao
}

// Leitura Formulario Fofoca
func FormularioFofoca(r *http.Request, apelidos map[string]crypto.Token, datahora time.Time) PostarFofoca {
	if r == nil {
		return PostarFofoca{}
	}
	if apelidos == nil {
		return PostarFofoca{}
	}
	acao := PostarFofoca{
		Acao:        "PostarFofoca",
		ID:          FormularioParaInt(r, "id"),
		CampoFofoca: r.FormValue("campoFofoca"),
		DataHora:    datahora,
	}
	return acao
}

// Leitura Formulario Ideia
func FormularioIdeia(r *http.Request, apelidos map[string]crypto.Token, datahora time.Time) PostarIdeia {
	if r == nil {
		return PostarIdeia{}
	}
	if apelidos == nil {
		return PostarIdeia{}
	}
	acao := PostarIdeia{
		Acao:       "PostarIdeia",
		ID:         FormularioParaInt(r, "id"),
		CampoIdeia: r.FormValue("campoIdeia"),
		DataHora:   datahora,
	}
	return acao
}

// Leitura Formulario Livro
func FormularioLivro(r *http.Request, apelidos map[string]crypto.Token, datahora time.Time, arquivo []byte, tipoArquivo string) PostarLivro {
	if r == nil {
		return PostarLivro{}
	}
	if apelidos == nil {
		return PostarLivro{}
	}
	acao := PostarLivro{
		Acao:         "PostarLivro",
		ID:           FormularioParaInt(r, "id"),
		TipoArquivo:  tipoArquivo,
		ArquivoLivro: arquivo,
		DataHora:     datahora,
	}
	return acao
}

// Leitura Formulario Meme
func FormularioMeme(r *http.Request, apelidos map[string]crypto.Token, datahora time.Time, arquivo []byte, tipoArquivo string) PostarMeme {
	if r == nil {
		return PostarMeme{}
	}
	if apelidos == nil {
		return PostarMeme{}
	}
	acao := PostarMeme{
		Acao:        "PostarMeme",
		ID:          FormularioParaInt(r, "id"),
		TipoArquivo: tipoArquivo,
		ArquivoMeme: arquivo,
		DataHora:    datahora,
	}
	return acao
}

// Leitura Formulario Musica
func FormularioMusica(r *http.Request, apelidos map[string]crypto.Token, datahora time.Time) PostarMusica {
	if r == nil {
		return PostarMusica{}
	}
	if apelidos == nil {
		return PostarMusica{}
	}
	acao := PostarMusica{
		Acao:        "PostarMusica",
		ID:          FormularioParaInt(r, "id"),
		CampoMusica: r.FormValue("campoMusica"),
		DataHora:    datahora,
	}
	return acao
}

// Retorna o tipo de arquivo
func TipoArquivo(nomeArquivo string) string {
	pos := strings.LastIndex(nomeArquivo, ".")
	if pos < len(nomeArquivo) && pos > 0 {
		return nomeArquivo[pos+1:]
	}
	return ""
}
*/
