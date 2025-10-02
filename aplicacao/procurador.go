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
	pk        crypto.PrivateKey
	signin    *GerenciadorSignin
	Token     crypto.Token
	carteira  crypto.PrivateKey
	gateway   chan []byte
	estado    *estado.Estado
	templates *template.Template
	indexer   *Index
	// senhaEmail   string
	genesisTime  time.Time
	ephemeralprv crypto.PrivateKey
	ephemeralpub crypto.Token
	nomeMucua    string
	hostname     string
	session      *configuracoes.CookieStore
	safe         *safe.Safe
	convidar     map[crypto.Hash]struct{} // map of invite user hash to token
}

func (a *ProcuradorGeral) Arroba(r *http.Request) string {
	autor := a.Autor(r)
	arroba := a.estado.HashTokenPraArrobas[crypto.HashToken(autor)]
	return arroba
}

func (a *ProcuradorGeral) Send(all []acoes.Acao, author crypto.Token) {
	for _, action := range all {
		dressed := a.DressAction(action, author)
		fmt.Println("Dressed action:", string(dressed))
		a.gateway <- append([]byte{messages.MsgAction}, dressed...)
	}
}

func (a *ProcuradorGeral) DressAction(action acoes.Acao, author crypto.Token) []byte {
	bytes := MegaParaBreeze(action.Serializa(), a.estado.Epoca)
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

func MegaParaBreeze(action []byte, epoch uint64) []byte {
	if action == nil {
		log.Print("PANIC BUG: MegaParaBreeze chamado com acao nula ")
		return nil
	}
	bytes := []byte{0, breeze.IVoid}                     // Breeze Void instruction version 0
	util.PutUint64(epoch, &bytes)                        // epoch (mega)
	bytes = append(bytes, 1, 1, 0, 0, attorney.VoidType) // mega protocol code + axe Void instruction code
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

func (a *ProcuradorGeral) DefineEpoca(epoca uint64) {
	a.estado.Epoca = epoca
}

func (a *ProcuradorGeral) IncorporarProcuracao(arroba string, procuracao *attorney.GrantPowerOfAttorney) {
	if procuracao != nil {
		a.signin.DarProcuracao(procuracao.Author, arroba, string(procuracao.Fingerprint))
	}
}

func (a *ProcuradorGeral) IncorporarRevogacao(arroba string, revogacao *attorney.RevokePowerOfAttorney) {
	if revogacao != nil {
		a.signin.RevogarProcuracao(revogacao.Author, arroba, revogacao.Signature.String())
	}
}

func (a *ProcuradorGeral) Incorporar(acao []byte) {
	if err := a.estado.Acao(acao); err != nil {
		fmt.Println("Erro ao incorporar acao:", err, acao)
	}
}
