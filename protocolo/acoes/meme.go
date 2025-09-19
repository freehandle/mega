package acoes

import (
	"time"

	"github.com/freehandle/breeze/crypto"
	"github.com/freehandle/breeze/util"
)

// estrutura da acao postar
type PostarMeme struct {
	Epoca       uint64
	Autor       crypto.Token
	TipoArquivo string
	Conteudo    []byte
	Data        time.Time
}

// faz o hash da instrucao
func (p *PostarMeme) FazHash() crypto.Hash {
	return crypto.Hasher([]byte(p.Serializa()))
}

// traz a autoria da instrucao
func (p *PostarMeme) Autoria() crypto.Token {
	return p.Autor
}

// serializa a instrucao
func (p *PostarMeme) Serializa() []byte {
	bytes := make([]byte, 0)
	util.PutUint64(p.Epoca, &bytes)
	util.PutToken(p.Autor, &bytes)
	util.PutByte(APostarMeme, &bytes)
	util.PutString(p.TipoArquivo, &bytes)
	util.PutByteArray(p.Conteudo, &bytes)
	util.PutTime(p.Data, &bytes)
	return bytes
}

// le a instrucao a partir da serializacao
func LeMeme(postmeme []byte) *PostarMeme {
	acao := PostarMeme{}
	posicao := 0
	acao.Epoca, posicao = util.ParseUint64(postmeme, posicao)
	acao.Autor, posicao = util.ParseToken(postmeme, posicao)
	// se o byte correspondente ao tipo de acao nao for o esperado, retorna nulo
	if postmeme[posicao] != APostarMeme {
		return nil
	}
	posicao += 1
	acao.TipoArquivo, posicao = util.ParseString(postmeme, posicao)
	acao.Conteudo, posicao = util.ParseByteArray(postmeme, posicao)
	acao.Data, posicao = util.ParseTime(postmeme, posicao)
	if posicao != len(postmeme) {
		return nil
	}
	return &acao
}
