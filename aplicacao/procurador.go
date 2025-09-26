package aplicacao

import (
	"fmt"
	"html/template"
	"log"
	"mega/aplicacao/configuracoes"
	"mega/protocolo/acoes"
	"mega/protocolo/estado"
	"net/http"
	"time"

	"github.com/freehandle/breeze/consensus/messages"
	"github.com/freehandle/breeze/crypto"
	breeze "github.com/freehandle/breeze/protocol/actions"
	"github.com/freehandle/breeze/util"
	"github.com/freehandle/handles/attorney"
	"github.com/freehandle/safe"
)

const cookieName = "sessaoMEGA"

type ProcuradorGeral struct {
	pk            crypto.PrivateKey
	signin        *GerenciadorSignin
	chavePrivada  crypto.PrivateKey
	Token         crypto.Token
	carteira      crypto.PrivateKey
	gateway       chan []byte
	estado        *estado.Estado
	templates     *template.Template
	indexer       *Index
	emailPassword string
	genesisTime   time.Time
	ephemeralprv  crypto.PrivateKey
	ephemeralpub  crypto.Token
	serverName    string
	hostname      string
	session       *configuracoes.CookieStore
	safe          *safe.Safe
	convidar      map[crypto.Hash]struct{} // map of invite user hash to token
}

func (a *ProcuradorGeral) Send(all []acoes.Acao, author crypto.Token) {
	for _, action := range all {
		dressed := a.DressAction(action, author)
		fmt.Println("Dressed action:", dressed)
		a.gateway <- append([]byte{messages.MsgAction}, dressed...)
	}
}

func (a *ProcuradorGeral) DressAction(action acoes.Acao, author crypto.Token) []byte {
	bytes := MegaToBreeze(action.Serializa(), a.estado.Epoca)
	if bytes == nil {
		return nil
	}
	for n := 0; n < crypto.TokenSize; n++ {
		bytes[15+n] = author[n]
	}

	// put attorney
	util.PutToken(a.pk.PublicKey(), &bytes)
	signature := a.pk.Sign(bytes)
	util.PutSignature(signature, &bytes)

	// put zero token wallet
	util.PutToken(a.pk.PublicKey(), &bytes)
	util.PutUint64(0, &bytes) // zero fee
	signature = a.pk.Sign(bytes)
	util.PutSignature(signature, &bytes)
	return bytes
}

func MegaToBreeze(action []byte, epoch uint64) []byte {
	if action == nil {
		log.Print("PANIC BUG: SynergyToBreeze called with nil action ")
		return nil
	}
	bytes := []byte{0, breeze.IVoid}                     // Breeze Void instruction version 0
	util.PutUint64(epoch, &bytes)                        // epoch (synergy)
	bytes = append(bytes, 1, 1, 0, 0, attorney.VoidType) // synergy protocol code + axe Void instruction code
	bytes = append(bytes, action[8:]...)                 //
	return bytes
}

func (a *ProcuradorGeral) Autor(r *http.Request) crypto.Token {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return crypto.ZeroToken
	}
	if token, ok := a.session.Get(cookie.Value); ok {
		return token
	}
	return crypto.ZeroToken
}
