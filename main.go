package main

import (
	"context"
	"log"
	"mega/aplicacao"
	"mega/aplicacao/configuracoes"
	"mega/protocolo/estado"
	"time"

	"github.com/freehandle/breeze/crypto"
	"github.com/freehandle/breeze/middleware/social"
	"github.com/freehandle/breeze/util"
	"github.com/freehandle/handles/attorney"
	"github.com/freehandle/safe"
)

const (
	handlesProtocolCode = 1
	breezeProtocolCode  = 0
	notarypath          = ""
	blocksPath          = ""
	blocksName          = "chain"
)

type ByArraySender chan []byte

func (b ByArraySender) Send(data []byte) error {
	b <- data
	return nil
}

func launchLocalChain(ctx context.Context, listeners []chan []byte, receiver chan []byte) error {
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

func launchMegaServer(gateway chan []byte, receive chan []byte, synergyPass, emailPass string, vault *config.SecretsVault, safe *safe.Safe) {
	indexador := aplicacao.NovoIndex()
	genesis := estado.EstadoInicial()
	indexador.InicializaEstado(genesis)

	attorneySecret := vault.PK
	cookieStore := configuracoes.OpenCokieStore("cookies.dat", genesis)
	passwordManager := configuracoes.NewFilePasswordManager("passwords.dat")
	config := aplicacao.ServerConfig{
		Vault:       vault,
		Attorney:    attorneySecret.PublicKey(),
		Ephemeral:   attorneySecret.PublicKey(),
		Passwords:   passwordManager,
		CookieStore: cookieStore,
		Indexer:     indexador,
		Gateway:     gateway,
		State:       genesis,
		GenesisTime: genesis.HoraGenesis,
		Hostname:    "localhost:3000",
		Mail:        &aplicacao.SMTPGmail{From: "freemyhandle@gmail.com", Password: emailPass},
		Port:        3000,
		Safe:        safe,
		//ServerName:    "/synergy",
	}
	attorney, finalize := aplicacao.NovoServidorProcuradorGeral(config)
	if attorney == nil {
		err := <-finalize
		log.Fatalf("error creating attorney: %v\n", err)
		return
	}
	handles := &network.HandlesDB{
		TokenToHandle: make(map[crypto.Token]network.UserInfo),
		HandleToToken: make(map[string]crypto.Token),
		Attorneys:     make(map[crypto.Token]struct{}),
		SynergyApp:    attorneySecret.PublicKey(),
	}
	genesis.Axe = handles
	signal := network.ByteArrayToSignal(receive)
	network.NewSynergyNode(handles, attorney, signal)
}

func main() {
	//
	return
}
