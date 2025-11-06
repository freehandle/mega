package acoes

import (
	"time"

	"github.com/freehandle/breeze/crypto"
	"github.com/freehandle/breeze/util"
)

// estrutura da acao aderir ao protocolo MIGA
type AderirMIGA struct {
	Epoca uint64
	Autor crypto.Token
	Data  time.Time
}

// faz o hash da instrucao
func (p *AderirMIGA) FazHash() crypto.Hash {
	return crypto.Hasher([]byte(p.Serializa()))
}

// traz a autoria da instrucao
func (p *AderirMIGA) Autoria() crypto.Token {
	return p.Autor
}

// serializa a instrucao
func (p *AderirMIGA) Serializa() []byte {
	bytes := make([]byte, 0)
	util.PutUint64(p.Epoca, &bytes)
	util.PutToken(p.Autor, &bytes)
	util.PutByte(AAderirMIGA, &bytes)
	util.PutTime(p.Data, &bytes)
	return bytes
}

// le a instrucao a partir da serializacao
func LeAderir(postideia []byte) *AderirMIGA {
	acao := AderirMIGA{}
	posicao := 0
	acao.Epoca, posicao = util.ParseUint64(postideia, posicao)
	acao.Autor, posicao = util.ParseToken(postideia, posicao)
	// se o byte correspondente ao tipo de acao nao for o esperado, retorna nulo
	if postideia[posicao] != AAderirMIGA {
		return nil
	}
	posicao += 1
	acao.Data, posicao = util.ParseTime(postideia, posicao)
	if posicao != len(postideia) {
		return nil
	}
	return &acao
}
