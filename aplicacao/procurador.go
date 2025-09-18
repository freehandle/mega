package aplicacao

import (
	"html/template"
	"mega/protocolo/estado"
	"time"

	"github.com/freehandle/breeze/crypto"
	"github.com/freehandle/breeze/middleware/blockdb/index"
	"github.com/freehandle/handles/attorney"
	"github.com/freehandle/safe"
)

type ProcuradorGeral struct {
	signin        *GerenciadorSignin
	chavePrivada  crypto.PrivateKey
	Token         crypto.Token
	carteira      crypto.PrivateKey
	gateway       chan []byte
	estado        *estado.Estado
	templates     *template.Template
	indexer       *index.Index
	emailPassword string
	genesisTime   time.Time
	ephemeralprv  crypto.PrivateKey
	ephemeralpub  crypto.Token
	serverName    string
	hostname      string
}

type Signerin struct {
	Handle      string
	Email       string
	TimeStamp   time.Time
	FingerPrint string
}

type GerenciadorSignin struct {
	safe    *safe.Safe // for optional direct onboarding
	pending []*Signerin
	// passwords     PasswordManager
	AttorneyToken crypto.Token
	// mail          Mailer
	Attorney *ProcuradorGeral
	Granted  map[string]crypto.Token
}

func (a *ProcuradorGeral) IncorporarProcuracao(arroba string, procuracao *attorney.GrantPowerOfAttorney) {
	if procuracao != nil {
		a.signin.GrantAttorney(procuracao.Author, arroba, string(procuracao.Fingerprint))
	}
}

func (a *ProcuradorGeral) RegisterAxeDataBase(axe state.HandleProvider) {
	a.state.Axe = axe
}
