package app

import (
	"log"

	"github.com/freehandle/breeze/crypto"
	breeze "github.com/freehandle/breeze/protocol/actions"
	"github.com/freehandle/breeze/socket"
	"github.com/freehandle/breeze/util"
	"github.com/freehandle/handles/attorney"
	"github.com/freehandle/mega/protocolo/acoes"
)

type PorteiraLocal chan []byte

type PorteiraRemota struct {
	Conexao *socket.SignedConnection
}

func (p *PorteiraRemota) Send(data []byte) error {
	data = append([]byte{0}, data...)
	return p.Conexao.Send(data)
}

func PorteiraInternet(conexao *socket.SignedConnection, credenciais crypto.PrivateKey) *Porteira {
	porteiraRemota := PorteiraRemota{Conexao: conexao}
	return &Porteira{portao: &porteiraRemota, credenciais: credenciais}
}

func (p PorteiraLocal) Send(data []byte) error {
	p <- data
	return nil
}

func (p PorteiraLocal) Close() error {
	close(p)
	return nil
}

type Portao interface {
	Send([]byte) error
}

func PorteiraDeCanal(canal chan []byte, credenciais crypto.PrivateKey) *Porteira {
	return &Porteira{
		portao:      PorteiraLocal(canal),
		credenciais: credenciais,
	}
}

type Porteira struct {
	portao      Portao
	credenciais crypto.PrivateKey
}

func MegaParaBreeze(action []byte, epoch uint64) []byte {
	if action == nil {
		log.Print("PANIC BUG: MegaParaBreeze chamado com acao nula ")
		return nil
	}
	bytes := []byte{0, breeze.IVoid}                     // Breeze Void instruction version 0
	util.PutUint64(epoch+256, &bytes)                    // epoch (mega)
	bytes = append(bytes, 1, 1, 0, 0, attorney.VoidType) // mega protocol code + palcos Void instruction code
	bytes = append(bytes, action[8:]...)                 //
	return bytes
}

func (p *Porteira) Encaminha(all []acoes.Acao, autoria crypto.Token, epoca uint64) {
	for _, action := range all {
		dressed := p.TravesteAcao(action, autoria, epoca)
		p.portao.Send(dressed)
		// gambiarra, depois usar o de baixo
		//p.portao.Send(append([]byte{messages.MsgAction}, dressed...))
	}
}

func (p *Porteira) TravesteAcao(action acoes.Acao, autoria crypto.Token, epoca uint64) []byte {
	bytes := MegaParaBreeze(action.Serializa(), epoca)
	if bytes == nil {
		return nil
	}
	// for n := 0; n < crypto.TokenSize; n++ {
	// 	bytes[15+n] = autoria[n]
	// }
	// put attorney
	util.PutToken(p.credenciais.PublicKey(), &bytes)
	signature := p.credenciais.Sign(bytes)
	util.PutSignature(signature, &bytes)
	util.PutToken(p.credenciais.PublicKey(), &bytes)
	util.PutUint64(0, &bytes)
	signature = p.credenciais.Sign(bytes)
	util.PutSignature(signature, &bytes)
	return bytes
}
