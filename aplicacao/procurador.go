package aplicacao

import (
	"text/template"
	"time"

	"github.com/freehandle/breeze/crypto"
	"github.com/freehandle/breeze/middleware/blockdb/index"
)

type AttorneyGeneral struct {
	chavePrivada  crypto.PrivateKey
	Token         crypto.Token
	carteira      crypto.PrivateKey
	gateway       chan []byte
	estado        *protocolo.Estado
	templates     *template.Template
	indexer       *index.Index
	emailPassword string
	genesisTime   time.Time
	ephemeralprv  crypto.PrivateKey
	ephemeralpub  crypto.Token
	serverName    string
	hostname      string
}
