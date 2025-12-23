package main

import (
	"context"
	"log"
	"time"

	"github.com/freehandle/breeze/crypto"
	"github.com/freehandle/breeze/middleware/simple"
	"github.com/freehandle/mega/app"
	"github.com/freehandle/mega/indice"
	"github.com/freehandle/mega/protocolo/estado"
)

func main() {
	pk := crypto.PrivateKeyFromString("e18e6528bd958000e51553f1828456c96509a3daa595421e24890d3153962297bb46f0c6a41ffc8ca179f3429d2584f103f66e540e21a197a45295ca8aa045de")
	token := pk.PublicKey()
	breezeToken := crypto.TokenFromString("91ad274d06c4be307a332a0e59449ad25ae2c65e4ad5a8f0af87067ac2fc3a54")

	ctx := context.Background()
	aplicacao := app.NovaAplicacaoVazia()
	novidades := simple.DissociateActions(ctx, simple.NewBlockReader(ctx, "/home/lienko/setembro/handles/cmd/proxy-handles", "blocos", time.Second))
	sender, err := simple.Gateway(ctx, 7000, breezeToken, pk)
	if err != nil {
		log.Fatalf("error creating gateway: %v", err)
	}
	aplicacao.Credenciais = pk
	aplicacao.Token = token
	aplicacao.Novidades = novidades
	aplicacao.Estado = estado.Genesis(0)
	aplicacao.Indice = indice.NovoIndice()
	aplicacao.Gateway = app.PorteiraDeCanal(sender, pk)
	aplicacao.NomeMucua = ""
	aplicacao.Hostname = ""
	aplicacao.Gerente, err = app.ContrataGerente(aplicacao, ".", "", "", pk)
	if err != nil {
		panic(err)
	}
	fim := make(chan error, 1)
	app.NovaMucua(ctx, aplicacao, 8070, "./app/static/", "localhost")
	err = <-fim
}

/*const (
	notarypath = ""
)

type ByArraySender chan []byte

func (b ByArraySender) Send(data []byte) error {
	b <- data
	return nil
}

func iniciaChainLocal(ctx context.Context, listeners []chan []byte, receiver chan []byte) error {
	genesis := attorney.NewGenesisState(notarypath)
	IO, err := util.OpenMultiFileStore(".", "blocos")
	if err != nil {
		return err
	}
	defer func() {
		IO.Close()
		log.Println("blockchain IO closed")
	}()

	chain := &social.LocalBlockChain[*attorney.Mutations, *attorney.MutatingState]{
		Interval:  time.Second,
		Listeners: listeners,
		Receiver:  receiver,
		IO:        IO,
	}
	if err = chain.LoadState(genesis, IO, listeners); err != nil {
		return err
	}
	return <-chain.Start(ctx)
}

func launchMegaServer(gateway chan []byte, receive chan []byte, megaPass, emailPass string, vault *configuracoes.SecretsVault, safe *safe.Safe) {
	indexador := aplicacao.NovoIndex()
	genesis := estado.EstadoInicial()
	indexador.InicializaEstado(genesis)

	attorneySecret := vault.PK
	cookieStore := configuracoes.OpenCokieStore("cookies.dat", genesis)
	passwordManager := configuracoes.NewFilePasswordManager("senhas.dat")
	config := aplicacao.ConfiguracaoMucua{
		Vault:       vault,
		Procurador:  attorneySecret.PublicKey(),
		Ephemeral:   attorneySecret.PublicKey(),
		Senhas:      passwordManager,
		CookieStore: cookieStore,
		Indexer:     indexador,
		Gateway:     gateway,
		State:       genesis,
		GenesisTime: genesis.HoraGenesis,
		Hostname:    "localhost:3000",
		Mail:        &aplicacao.SMTPGmail{From: "freemyhandle@gmail.com", Password: emailPass},
		Port:        3000,
		Safe:        safe,
		NomeMucua:   "/mega",
	}
	attorney, finalize := aplicacao.NovaMucuaProcuradorGeral(config)
	if attorney == nil {
		err := <-finalize
		log.Fatalf("error creating attorney: %v\n", err)
		return
	}
	handles := &rede.HandlesDB{
		TokenToHandle: make(map[crypto.Token]rede.UserInfo),
		HandleToToken: make(map[string]crypto.Token),
		Attorneys:     make(map[crypto.Token]struct{}),
		AppMEGA:       attorneySecret.PublicKey(),
	}
	genesis.Apelidos = handles
	signal := rede.ByteArrayToSignal(receive)
	rede.NewMEGANode(handles, attorney, signal)
}

func main() {
	envs := os.Environ()
	var emailSenha string
	var cofreSenha string
	for _, env := range envs {
		if strings.HasPrefix(env, "SENHA_EMAIL=") {
			emailSenha, _ = strings.CutPrefix(env, "SENHA_EMAIL=")
		} else if strings.HasPrefix(env, "SENHA_COFRE=") {
			cofreSenha, _ = strings.CutPrefix(env, "SENHA_COFRE=")
		}
	}

	monitorCofre := make(chan []byte)
	monitorMEGA := make(chan []byte)
	fornecedor := make(chan []byte)

	ctxBack := context.Background()
	ctx, cancel := context.WithCancel(ctxBack)

	go iniciaChainLocal(ctx, []chan []byte{monitorMEGA, monitorCofre}, fornecedor)

	vault, err := configuracoes.OpenVaultFromPassword([]byte(cofreSenha), "megavault.dat")
	if err != nil {
		log.Fatalf("erro ao abrir o vault: %v \n", err)
		return
	}
	vault.Secrets[vault.PK.PublicKey()] = vault.PK

	cfg := safe.SafeConfig{
		Credentials: vault.PK,
		HtmlPath:    "../safe/",
		Path:        ".",
		Port:        7000,
		ServerName:  "/safe",
	}
	errSignal, safe := safe.NewLocalServer(ctx, cfg, cofreSenha, ByArraySender(fornecedor), monitorCofre)

	go func() {
		err := <-errSignal
		log.Printf("error creating safe server: %v", err)
		cancel()
	}()

	go launchMegaServer(fornecedor, monitorMEGA, cofreSenha, emailSenha, vault, safe)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	s := <-c
	fmt.Println("Got signal:", s)
	cancel()
	time.Sleep(5 * time.Second)
}
*/
