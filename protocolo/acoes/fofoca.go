package acoes

import (
	"time"
	"unicode/utf8"

	"github.com/freehandle/breeze/crypto"
	"github.com/freehandle/breeze/util"
)

// estrutura da acao postar ideia
type PostarFofoca struct {
	Epoca    uint64
	Autor    crypto.Token
	Conteudo string // texto 2 paginas
	Data     time.Time
}

func (p *PostarFofoca) ValidarFormato() bool {
	return utf8.RuneCountInString(p.Conteudo) >= 100 && utf8.RuneCountInString(p.Conteudo) <= 5000
}

// faz o hash da instrucao
func (p *PostarFofoca) FazHash() crypto.Hash {
	return crypto.Hasher([]byte(p.Serializa()))
}

// traz a autoria da instrucao
func (p *PostarFofoca) Autoria() crypto.Token {
	return p.Autor
}

// serializa a instrucao
func (p *PostarFofoca) Serializa() []byte {
	bytes := make([]byte, 0)
	util.PutUint64(p.Epoca, &bytes)
	util.PutToken(p.Autor, &bytes)
	util.PutByte(APostarFofoca, &bytes)
	util.PutString(p.Conteudo, &bytes)
	util.PutTime(p.Data, &bytes)
	return bytes
}

// le a instrucao a partir da serializacao
func LeFofoca(postideia []byte) *PostarFofoca {
	acao := PostarFofoca{}
	posicao := 0
	acao.Epoca, posicao = util.ParseUint64(postideia, posicao)
	acao.Autor, posicao = util.ParseToken(postideia, posicao)
	// se o byte correspondente ao tipo de acao nao for o esperado, retorna nulo
	if postideia[posicao] != APostarFofoca {
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
