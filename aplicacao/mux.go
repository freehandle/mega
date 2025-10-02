package aplicacao

import (
	"fmt"
	"html/template"
	"log"
	"mega/aplicacao/configuracoes"
	"mega/protocolo/estado"
	"net/http"
	"time"

	"github.com/freehandle/breeze/crypto"
	"github.com/freehandle/safe"
)

var arquivosTemplate []string = []string{
	"main", "verjornal", "criartxt",
}

type ConfiguracaoMucua struct {
	Vault       *configuracoes.SecretsVault
	Procurador  crypto.Token
	Ephemeral   crypto.Token
	Senhas      configuracoes.PasswordManager
	CookieStore *configuracoes.CookieStore
	Indexer     *Index
	Gateway     chan []byte
	State       *estado.Estado
	GenesisTime time.Time
	Mail        Mailer
	Port        int
	Path        string
	NomeMucua   string
	Hostname    string
	Safe        *safe.Safe
}

func NovaMucuaProcuradorGeral(cfg ConfiguracaoMucua) (*ProcuradorGeral, chan error) {
	finalize := make(chan error, 2)

	attorneySecret, ok := cfg.Vault.Secrets[cfg.Procurador]
	if !ok {
		finalize <- fmt.Errorf("attorney secret key not found in vault")
		return nil, finalize
	}
	ephemeralSecret, ok := cfg.Vault.Secrets[cfg.Ephemeral]
	if !ok {
		finalize <- fmt.Errorf("ephemeral secret key not found in vault")
		return nil, finalize
	}

	attorney := ProcuradorGeral{
		pk:           attorneySecret,
		Token:        cfg.Procurador,
		carteira:     attorneySecret,
		gateway:      cfg.Gateway,
		estado:       cfg.State,
		indexer:      cfg.Indexer,
		session:      cfg.CookieStore,
		genesisTime:  cfg.GenesisTime,
		ephemeralpub: cfg.Ephemeral,
		ephemeralprv: ephemeralSecret,
		nomeMucua:    cfg.NomeMucua,
		hostname:     cfg.Hostname,
		safe:         cfg.Safe,
		convidar:     make(map[crypto.Hash]struct{}),
	}
	if cfg.Path == "" {
		cfg.Path = "./"
	}
	templatesPath := fmt.Sprintf("%v/api/templates", cfg.Path)
	attorney.signin = NewSigninManager(cfg.Senhas, cfg.Mail, &attorney)
	attorney.templates = template.New("root")
	files := make([]string, len(arquivosTemplate))

	for n, file := range arquivosTemplate {
		files[n] = fmt.Sprintf("%v/%v.html", templatesPath, file)
	}
	t, err := template.ParseFiles(files...)
	if err != nil {
		log.Fatal(err)
	}
	attorney.templates = t

	staticPath := fmt.Sprintf("%v/api/static/", cfg.Path)
	go NovaMucua(&attorney, cfg.Port, staticPath, finalize, cfg.NomeMucua)

	return &attorney, finalize
}

func NovaMucua(procurador *ProcuradorGeral, port int, staticPath string, finalize chan error, servername string) {

	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir(staticPath))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	mux.HandleFunc("/", procurador.AgentePrincipal) // funcao que gera o template main
	mux.HandleFunc("/verjornal", procurador.AgenteVerJornal)
	mux.HandleFunc("/api", procurador.AgenteAPI)
	mux.HandleFunc("/uploadfile", procurador.OperadorUpload)
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%v", port),
		Handler:      mux,
		WriteTimeout: 2 * time.Second,
	}
	finalize <- srv.ListenAndServe()
}
