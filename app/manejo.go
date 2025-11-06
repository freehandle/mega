package app

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/freehandle/breeze/crypto"
)

type InformacaoCabecalho struct {
	Arroba          string
	Ativo           string
	Erro            string
	NomeMucua       string
	LinkSelecionada string
}

type ViewConvite struct {
	Cabecalho InformacaoCabecalho
	Seed      string
	// pra testar
	Nome  string
	Nome2 string
}

// Gerenciador do template principal da aplicacao
func (a *Aplicacao) ManejoPrincipal(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("X-Content-Type-Options", "nosniff")
	view := InformacaoCabecalho{
		Arroba:          a.Autor(r),
		NomeMucua:       a.NomeMucua,
		Ativo:           "",
		LinkSelecionada: "",
	}
	if err := a.templates.ExecuteTemplate(w, "main.html", view); err != nil {
		log.Println(err)
	}
}

func (a *Aplicacao) ManejoSignin(w http.ResponseWriter, r *http.Request) {
	hashEncoded := r.URL.Path
	hashEncoded = strings.Replace(hashEncoded, "/signin/", "", 1)
	hash := crypto.DecodeHash(hashEncoded)
	fmt.Println("oia", len(a.Convidar))
	if _, ok := a.Convidar[hash]; ok || len(a.Convidar) == 0 {
		view := ViewConvite{
			Cabecalho: InformacaoCabecalho{
				NomeMucua:       "",
				Ativo:           "",
				LinkSelecionada: "",
				Arroba:          "",
			},
			Seed:  hashEncoded,
			Nome:  "teste2",
			Nome2: "teste3",
		}
		fmt.Println("Seed:", hashEncoded)
		if err := a.templates.ExecuteTemplate(w, "signin.html", view); err != nil {
			log.Println(err)
		}
	} else {
		view := InformacaoCabecalho{
			Erro:      "convite inválido",
			NomeMucua: a.NomeMucua,
		}
		if err := a.templates.ExecuteTemplate(w, "login.html", view); err != nil {
			log.Println(err)
		}
	}
}

func (a *Aplicacao) ManejoNovoUsuario(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}
	arroba := r.FormValue("handle")
	email := r.FormValue("email")
	senha := r.FormValue("password")
	ok := a.Gerente.OnboardSigner(arroba, email, senha)
	aviso := InformacaoCabecalho{
		NomeMucua: a.NomeMucua,
	}
	if !ok {
		aviso.Erro = "Confira seu email para ativar sua conta ou tente outro arroba."
	}
	if err := a.templates.ExecuteTemplate(w, "login.html", aviso); err != nil {
		log.Println(err)
	}
	return
}

func (a *Aplicacao) ManejoCredenciais(w http.ResponseWriter, r *http.Request) {
	cookie, arroba, err := a.Gerente.CredentialsHandler(r)
	fmt.Println(cookie)
	fmt.Println(arroba)
	if err != nil {
		aviso := InformacaoCabecalho{
			NomeMucua: a.NomeMucua,
			Erro:      err.Error(),
		}
		if err := a.templates.ExecuteTemplate(w, "login.html", aviso); err != nil {
			log.Println(err)
		}
		return
	}
	http.SetCookie(w, cookie)
	aviso := InformacaoCabecalho{
		NomeMucua: a.NomeMucua,
		Arroba:    arroba,
	}
	if err := a.templates.ExecuteTemplate(w, "main.html", aviso); err != nil {
		log.Println(err)
	}
	return
}
