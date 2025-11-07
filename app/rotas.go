package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/freehandle/breeze/crypto"
	"github.com/freehandle/mega/aplicacao/configuracoes"
	"github.com/freehandle/mega/protocolo/estado"
	"github.com/freehandle/safe"
)

type ConfiguracaoMucua struct {
	Vault      *configuracoes.SecretsVault
	Procurador crypto.Token
	//Ephemeral   crypto.Token
	//Senhas      configuracoes.PasswordManager
	//CookieStore *configuracoes.CookieStore
	//Indexer     *Index
	Gateway     chan []byte
	State       *estado.Estado
	GenesisTime time.Time
	//Mail        Mailer
	Port      int
	Path      string
	NomeMucua string
	Hostname  string
	Safe      *safe.Safe
}

func NovaMucua(ctx context.Context, app *Aplicacao, port int, staticPath string, servername string) {
	go app.Rodar(ctx)
	mux := http.NewServeMux()
	fmt.Println("Static path:", staticPath)
	fs := http.FileServer(http.Dir(staticPath))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("/", app.ManejoPrincipal) // funcao que gera o template main
	//mux.HandleFunc("/verjornal", procurador.AgenteVerJornal)
	mux.HandleFunc("/signin", app.ManejoSignin)
	mux.HandleFunc("/novousuario", app.ManejoNovoUsuario)
	mux.HandleFunc("/credenciais", app.ManejoCredenciais)
	mux.HandleFunc("/publicar", app.ManejoInterfacePublicar)
	mux.HandleFunc("/publica", app.ManejoPublica)

	//mux.HandleFunc("/uploadfile", procurador.OperadorUpload)
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%v", port),
		Handler:      mux,
		WriteTimeout: 2 * time.Second,
	}
	srv.ListenAndServe()
}
