package main

import (
	"log"
	"time"

	"github.com/freehandle/iu/config"
	"github.com/freehandle/mega/app"
	"github.com/freehandle/mega/indice"
	"github.com/freehandle/mega/protocolo/estado"
)

func main() {

	aplicacao := app.NovaAplicacaoVazia()
	recursos, erro := config.LoadConfig("mega.json")
	if erro != nil {
		log.Fatal(erro)
	}

	aplicacao.Credenciais = recursos.Secret
	aplicacao.Token = recursos.Secret.PublicKey()
	aplicacao.Novidades = recursos.Blocks
	aplicacao.Estado = estado.Genesis(0)
	aplicacao.Indice = indice.NovoIndice()
	aplicacao.GenesisTime = recursos.GenesisTime
	aplicacao.Intervalo = time.Second
	aplicacao.Gateway = app.PorteiraDeCanal(recursos.Gateway, recursos.Secret)
	aplicacao.NomeMucua = recursos.Config.Mucua
	aplicacao.Hostname = recursos.Config.HostName
	aplicacao.CaminhoArquivos = recursos.Config.UserFilesPath
	aplicacao.Gerente = recursos.Manager
	fim := make(chan error, 1)
	app.NovaMucua(recursos.Context, aplicacao, recursos.Config.AppPort, "./app", "localhost")
	erro = <-fim
	if erro != nil {
		log.Fatal(erro)
	}
}
