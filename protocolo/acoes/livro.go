package acoes

import (
	"time"

	"github.com/freehandle/breeze/crypto"
	"github.com/freehandle/breeze/util"
)

// estrutura da acao postar
type PostarLivro struct {
	Epoca    uint64
	Autor    crypto.Token
	Conteudo []byte // imagem (ver em mega)
	Data     time.Time
}

// faz o hash da instrucao
func (p *PostarLivro) FazHash() crypto.Hash {
	return crypto.Hasher([]byte(p.Serializa()))
}

// traz a autoria da instrucao
func (p *PostarLivro) Autoria() crypto.Token {
	return p.Autor
}

// serializa a instrucao
func (p *PostarLivro) Serializa() []byte {
	bytes := make([]byte, 0)
	util.PutUint64(p.Epoca, &bytes)
	util.PutToken(p.Autor, &bytes)
	util.PutByte(APostarLivro, &bytes)
	util.PutByteArray(p.Conteudo, &bytes)
	util.PutTime(p.Data, &bytes)
	return bytes
}

// le a instrucao a partir da serializacao
func LeLivro(postideia []byte) *PostarLivro {
	acao := PostarLivro{}
	posicao := 0
	acao.Epoca, posicao = util.ParseUint64(postideia, posicao)
	acao.Autor, posicao = util.ParseToken(postideia, posicao)
	// se o byte correspondente ao tipo de acao nao for o esperado, retorna nulo
	if postideia[posicao] != APostarLivro {
		return nil
	}
	posicao += 1
	acao.Conteudo, posicao = util.ParseByteArray(postideia, posicao)
	acao.Data, posicao = util.ParseTime(postideia, posicao)
	if posicao != len(postideia) {
		return nil
	}
	return &acao
}
