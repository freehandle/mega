package aplicacao

import (
	"log"
	"net/http"
	"strings"
)

// Gerenciador do template principal da aplicacao
func (a *ProcuradorGeral) AgentePrincipal(w http.ResponseWriter, r *http.Request) {
	view := InformacaoCabecalho{
		ArrobaLogada: a.Arroba(r),
		NomeMucua:    a.nomeMucua,
		Ativo:        "",
	}
	if err := a.templates.ExecuteTemplate(w, "main.html", view); err != nil {
		log.Println(err)
	}
}

func (a *ProcuradorGeral) AgenteVerJornal(w http.ResponseWriter, r *http.Request) {
	cabecalho := InformacaoCabecalho{
		ArrobaLogada: a.Arroba(r),
		Ativo:        "Ver",
		NomeMucua:    a.nomeMucua,
	}
	arroba := r.URL.Path
	arroba = strings.Replace(arroba, "/jornal/", "", 1)
	view := JornalUsuarioDoEstado(a.estado, arroba)
	view.Cabecalho = cabecalho
	if err := a.templates.ExecuteTemplate(w, "verjornal.html", view); err != nil {
		log.Println(err)
	} else {
		return
	}
}

// Agentes para acoes de assinatura, login e reset de senha de usuarios

// Ingresso de usuarios na aplicacao
// func (a *ProcuradorGeral) AgenteNovoUsuario(w http.ResponseWriter, r *http.Request) {
// 	if err := r.ParseForm(); err != nil {
// 		fmt.Fprintf(w, "ParseForm() err: %v", err)
// 		return
// 	}
// 	email := r.FormValue("email")
// 	arroba := r.FormValue("handle")
// 	token := a.estado.Apelidos.Token(arroba)
// 	view := Mucua{
// 		Nome: a.nomeMucua,
// 		Cabecalho: InformacaoCabecalho{
// 			Erro: "você já é usuário: por favor, entre!",
// 		},
// 	}
// 	if token != nil {
// 		isMember := a.signin.senhas.Has(*token)
// 		if isMember {
// 			if err := a.templates.ExecuteTemplate(w, "login.html", view); err != nil {
// 				log.Println(err)
// 				http.Redirect(w, r, fmt.Sprintf("%v/login", a.nomeMucua), http.StatusSeeOther)
// 			}
// 			return
// 		}
// 	}
// 	a.signin.AddSigner(handle, email, token)
// 	http.Redirect(w, r, fmt.Sprintf("%v/login", a.serverName), http.StatusSeeOther)
// }

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
