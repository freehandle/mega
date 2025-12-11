package app

import (
	"context"
	"fmt"
	"html/template"
	"iu/auth"
	"log"
	"net/http"
	"time"

	"github.com/freehandle/breeze/crypto"
	"github.com/freehandle/breeze/util"
	"github.com/freehandle/handles/attorney"
	"github.com/freehandle/mega/indice"
	"github.com/freehandle/mega/protocolo/estado"
)

const appName = "MIGA"

var arquivosTemplate []string = []string{
	"login", "signin", "meujornal", "login", "novotexto",
}

type Aplicacao struct {
	Epoca       uint64
	Credenciais crypto.PrivateKey
	Token       crypto.Token
	Gateway     *Porteira
	Novidades   chan []byte
	Estado      *estado.Estado
	Indice      *indice.Indice
	templates   *template.Template
	GenesisTime time.Time
	NomeMucua   string
	Hostname    string
	Convidar    map[crypto.Hash]struct{} // map of invite user hash to token
	Gerente     *auth.SigninManager
}

func (p *Aplicacao) Rodar(ctx context.Context) {
	validador := p.Estado.Validator()
	for {
		select {
		case <-ctx.Done():
			log.Println("Aplicacao.Rodar: context done, exiting")
			return
		case novidade := <-p.Novidades:
			log.Printf("Aplicacao.Rodar: received action of size %d bytes\n", len(novidade))
			if len(novidade) == 0 {
				continue
			}
			if novidade[0] == 0 {
				if len(novidade) >= 9 {
					epoca, _ := util.ParseUint64(novidade[1:], 1)
					mutacoes := validador.Mutations()
					p.Estado.Incorporate(mutacoes)
					validador = p.Estado.Validator()
					validador.Mutacoes.Epoca = epoca
					p.Epoca = epoca
				} else {
					continue
				}
			} else {
				acao := novidade[1:]
				if tipoHandles := attorney.Kind(acao); tipoHandles == attorney.JoinNetworkType {
					if usuario := attorney.ParseJoinNetwork(acao); usuario != nil {
						p.Indice.IncorporaAutor(usuario.Handle, usuario.Author)
					}
				}
				if validador.Validate(acao) {
					p.Indice.IncorporaAcao(acao)
				}
			}
		}
	}
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
	if err == nil {
		if handle, ok := p.Gerente.Cookies.Get(cookie.Value); ok {
			return handle
		}
	}
	return ""
}
