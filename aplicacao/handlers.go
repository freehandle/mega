package aplicacao

import (
	"log"
	"net/http"
)

func (a *ProcuradorGeral) MainHandler(w http.ResponseWriter, r *http.Request) {
	view := Mucua{
		Cabecalho: InformacaoCabecalho{
			ArrobaUsuario: a.Arroba(r),
			NomeMucua:     a.nomeMucua,
		},
		Nome: a.nomeMucua,
	}
	if err := a.templates.ExecuteTemplate(w, "main.html", view); err != nil {
		log.Println(err)
	}
}

// func (a *ProcuradorGeral) ApiHandler(w http.ResponseWriter, r *http.Request) {
// 	if err := r.ParseForm(); err != nil {
// 		fmt.Fprintf(w, "ParseForm() err: %v", err)
// 		return
// 	}
// 	var actionArray []acoes.Acao
// 	var err error
// 	data := time.Now()
// 	autor := a.Autor(r)
// 	switch r.FormValue("action") {
// 	case "PostarCauso":
// 		actionArray, err = FormularioCauso(r, a.estado.ArrobasPraTokens, data).ParaAcao()
// 	case "PostarFofoca":
// 		actionArray, err = FormularioFofoca(r, a.estado.ArrobasPraTokens, data).ParaAcao()
// 	case "PostarIdeia":
// 		actionArray, err = FormularioIdeia(r, a.estado.ArrobasPraTokens, data).ParaAcao()
// 	// case "PostarLivro":
// 	// 	actionArray, err = FormularioLivro(r, a.estado.ArrobasPraTokens, data, arquivo, tipoArquivo).ParaAcao()
// 	// case "PostarMeme":
// 	// 	actionArray, err = FormularioMeme(r, a.estado.ArrobasPraTokens, data, arquivo, tipoArquivo).ParaAcao()
// 	case "PostarMusica":
// 		actionArray, err = FormularioMusica(r, a.estado.ArrobasPraTokens, data).ParaAcao()
// 	}
// 	if err == nil && len(actionArray) > 0 {
// 		a.Send(actionArray, autor)
// 	}
// 	redirect := fmt.Sprintf("%v/%v", a.serverName, r.FormValue("redirect"))
// 	http.Redirect(w, r, redirect, http.StatusSeeOther)
// }
