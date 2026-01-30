package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/freehandle/breeze/crypto"
	"github.com/freehandle/breeze/middleware/simple"
	"github.com/freehandle/mega/app"
	"github.com/freehandle/mega/indice"
	"github.com/freehandle/mega/protocolo/estado"
)

func main() {

	// Ajustando pra envio do email
	envs := os.Environ()
	var senhaEmail string
	caminhoBlocos := "/home/lienko/setembro/handles/cmd/proxy-handles" // caminho default da minha maquina

	for _, env := range envs {
		if strings.HasPrefix(env, "SENHA_EMAIL=") {
			senhaEmail, _ = strings.CutPrefix(env, "SENHA_EMAIL=")
		} else if strings.HasPrefix(env, "CAMINHO_BLOCOS=") {
			caminhoBlocos, _ = strings.CutPrefix(env, "CAMINHO_BLOCOS=")
		}
	}
	fmt.Println(senhaEmail)

	pk := crypto.PrivateKeyFromString("e18e6528bd958000e51553f1828456c96509a3daa595421e24890d3153962297bb46f0c6a41ffc8ca179f3429d2584f103f66e540e21a197a45295ca8aa045de")
	token := pk.PublicKey()
	breezeToken := crypto.TokenFromString("91ad274d06c4be307a332a0e59449ad25ae2c65e4ad5a8f0af87067ac2fc3a54")

	ctx := context.Background()
	aplicacao := app.NovaAplicacaoVazia()
	novidades := simple.DissociateActions(ctx, simple.NewBlockReader(ctx, caminhoBlocos, "blocos", time.Second))
	sender, err := simple.Gateway(ctx, 7000, breezeToken, pk)
	if err != nil {
		log.Fatalf("error creating gateway: %v", err)
	}
	// gatewayConn, err := socket.Dial("localhost", "localhost:7000", pk, breezeToken)
	// if err != nil {
	// 	panic(err)
	// }
	aplicacao.Credenciais = pk
	aplicacao.Token = token
	aplicacao.Novidades = novidades
	aplicacao.Estado = estado.Genesis(0)
	aplicacao.Indice = indice.NovoIndice()
	aplicacao.GenesisTime = time.Date(2025, time.December, 24, 15, 10, 10, 0, time.UTC) // colocado à mao
	aplicacao.Intervalo = time.Second
	aplicacao.Gateway = app.PorteiraDeCanal(sender, pk)
	// aplicacao.Gateway = app.PorteiraInternet(gatewayConn, pk)
	aplicacao.NomeMucua = ""
	aplicacao.Hostname = ""
	if senhaEmail == "" {
		aplicacao.Gerente, err = app.ContrataGerente(aplicacao, ".", "", "", pk)
	} else {
		aplicacao.Gerente, err = app.ContrataGerente(aplicacao, ".", senhaEmail, "arrobaslivres@gmail.com", pk)
	}
	if err != nil {
		panic(err)
	}
	// fmt.Println(aplicacao.DataDaEpoca(10))
	fim := make(chan error, 1)
	app.NovaMucua(ctx, aplicacao, 8070, "./app", "localhost")
	err = <-fim
}
