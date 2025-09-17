package aplicacao

import (
	"net/http"
	"strconv"
	"time"

	"github.com/freehandle/breeze/crypto"
)

// Conversao de tipos
func FormularioParaInt(r *http.Request, field string) int {
	if r == nil {
		return 0
	}
	value, _ := strconv.Atoi(r.FormValue(field))
	return value
}

// Leitura de formularios
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
		CampoCauso: r.FormValue("campocauso"),
		DataHora:   datahora,
	}
	return acao
}

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
		CampoFofoca: r.FormValue("campofofoca"),
		DataHora:    datahora,
	}
	return acao
}
