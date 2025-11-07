package acoes

import (
	"time"
	"unicode/utf8"

	"github.com/freehandle/breeze/crypto"
	"github.com/freehandle/breeze/util"
)

// estrutura da acao postar
type PostarCauso struct {
	Epoca    uint64
	Autor    crypto.Token
	Conteudo string // texto 1 paragrafo
	Data     time.Time
}

func (p *PostarCauso) ValidarFormato() bool {
	return utf8.RuneCountInString(p.Conteudo) >= 100 && utf8.RuneCountInString(p.Conteudo) <= 1600
}

// faz o hash da instrucao
func (p *PostarCauso) FazHash() crypto.Hash {
	return crypto.Hasher([]byte(p.Serializa()))
}

// traz a autoria da instrucao
func (p *PostarCauso) Autoria() crypto.Token {
	return p.Autor
}

// serializa a instrucao
func (p *PostarCauso) Serializa() []byte {
	bytes := make([]byte, 0)
	util.PutUint64(p.Epoca, &bytes)
	util.PutToken(p.Autor, &bytes)
	util.PutByte(APostarCauso, &bytes)
	util.PutString(p.Conteudo, &bytes)
	util.PutTime(p.Data, &bytes)
	return bytes
}

// le a instrucao a partir da serializacao
func LeCauso(postideia []byte) *PostarCauso {
	acao := PostarCauso{}
	posicao := 0
	acao.Epoca, posicao = util.ParseUint64(postideia, posicao)
	acao.Autor, posicao = util.ParseToken(postideia, posicao)
	// se o byte correspondente ao tipo de acao nao for o esperado, retorna nulo
	if postideia[posicao] != APostarCauso {
		return nil
	}
	posicao += 1
	acao.Conteudo, posicao = util.ParseString(postideia, posicao)
	acao.Data, posicao = util.ParseTime(postideia, posicao)
	if posicao != len(postideia) {
		return nil
	}
	return &acao
}
