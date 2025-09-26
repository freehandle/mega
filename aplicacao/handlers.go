package aplicacao

import (
	"fmt"
	"mega/protocolo/acoes"
	"net/http"
	"time"
)

func (a *ProcuradorGeral) ApiHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	var actionArray []acoes.Acao
	var err error
	data := time.Now()
	autor := a.Autor(r)
	switch r.FormValue("action") {
	case "PostarCauso":
		actionArray, err = FormularioCauso(r, a.estado.ArrobasPraTokens, data).ParaAcao()
	case "PostarFofoca":
		actionArray, err = FormularioFofoca(r, a.estado.ArrobasPraTokens, data).ParaAcao()
		// case "PostarIdeia":
		// 	actionArray, err = FormularioIdeia(r, a.estado.ArrobasPraTokens, data).ParaAcao()
		// case "PostarIdeia":
		// 	actionArray, err = FormularioIdeia(r, a.estado.ArrobasPraTokens, data).ParaAcao()
		// case "PostarLivro":
		// actionArray, err = FormularioLivro(r, a.estado.ArrobasPraTokens, data, arquivo, tipoArquivo).ParaAcao()
	}
	if err == nil && len(actionArray) > 0 {
		a.Send(actionArray, autor)
	}
	redirect := fmt.Sprintf("%v/%v", a.serverName, r.FormValue("redirect"))
	http.Redirect(w, r, redirect, http.StatusSeeOther)
}
