package app

import (
	"fmt"
	"html/template"
	"iu/auth"
	"log"
	"net/http"
	"time"

	"github.com/freehandle/breeze/crypto"
)

const appName = "MIGA"

var arquivosTemplate []string = []string{
	"main", "signin", "login",
}

type Aplicacao struct {
	Credenciais crypto.PrivateKey
	Token       crypto.Token
	Gateway     chan []byte
	//Estado      *estado.Estado
	templates *template.Template
	//indexer   *Index
	// senhaEmail   string
	GenesisTime time.Time
	NomeMucua   string
	Hostname    string
	Convidar    map[crypto.Hash]struct{} // map of invite user hash to token
	Gerente     *auth.SigninManager
}

func NovaAplicacaoVazia() *Aplicacao {
	files := make([]string, len(arquivosTemplate))
	for n, file := range arquivosTemplate {
		files[n] = fmt.Sprintf("%v/%v.html", "./aplicacao/templates", file)
	}
	t, err := template.ParseFiles(files...)
	if err != nil {
		log.Fatal(err)
	}

	return &Aplicacao{
		Convidar:    make(map[crypto.Hash]struct{}),
		templates:   t,
		GenesisTime: time.Now(),
	}
}

func (p *Aplicacao) Invite(handle string, token crypto.Token) error {
	return nil
}

func (p *Aplicacao) AppName() string {
	return appName
}

func (p *Aplicacao) AttorneyToken() crypto.Token {
	return p.Token
}

func (p *Aplicacao) Autor(r *http.Request) string {
	cookie, err := r.Cookie(appName)
	if err != nil {
		return ""
	}
	if handle, ok := p.Gerente.Cookies.Get(cookie.Value); ok {
		return handle
	}
	return ""
}
