package aplicacao

import (
	"fmt"
	"html/template"
	"log"
	"mime"
	"net/http"
	"time"

	"github.com/freehandle/mega/aplicacao/configuracoes"
	"github.com/freehandle/mega/protocolo/estado"

	"github.com/freehandle/breeze/crypto"
	"github.com/freehandle/safe"
)

var arquivosTemplate []string = []string{
	"main", "verjornal", "criartxt", "signin",
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
		cfg.Path = "."
	}
	templatesPath := fmt.Sprintf("%v/aplicacao/templates", cfg.Path)
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

	staticPath := fmt.Sprintf("%v/aplicacao/static/", cfg.Path)
	go NovaMucua(&attorney, cfg.Port, staticPath, finalize, cfg.NomeMucua)

	return &attorney, finalize
}

func NovaMucua(procurador *ProcuradorGeral, port int, staticPath string, finalize chan error, servername string) {

	// mime.AddExtensionType(".js", "application/javascript")
	mime.AddExtensionType(".css", "text/css; charset=utf-8")

	mux := http.NewServeMux()
	fmt.Println("Static path:", staticPath)
	fs := http.FileServer(http.Dir(staticPath))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// csss := http.FileServer(http.Dir(staticPath))
	mux.HandleFunc("/", procurador.AgentePrincipal) // funcao que gera o template main
	mux.HandleFunc("/verjornal", procurador.AgenteVerJornal)
	mux.HandleFunc("/signin", procurador.AgenteSignin)

	mux.HandleFunc("/api", procurador.AgenteAPI)
	mux.HandleFunc("/uploadfile", procurador.OperadorUpload)
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%v", port),
		Handler:      mux,
		WriteTimeout: 2 * time.Second,
	}
	finalize <- srv.ListenAndServe()
}
