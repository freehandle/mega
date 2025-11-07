package acoes

import (
	"slices"
	"strings"
	"time"

	"github.com/freehandle/breeze/crypto"
	"github.com/freehandle/breeze/util"
)

// estrutura da acao postar
type PostarLivro struct {
	Epoca       uint64
	Autor       crypto.Token
	TipoArquivo string
	Conteudo    crypto.Hash
	Data        time.Time
}

func (p *PostarLivro) ValidarFormato() bool {
	tipomin := strings.ToLower(p.TipoArquivo)
	return slices.Contains(TiposImagens, tipomin)
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
	util.PutString(p.TipoArquivo, &bytes)
	util.PutHash(p.Conteudo, &bytes)
	util.PutTime(p.Data, &bytes)
	return bytes
}

// le a instrucao a partir da serializacao
func LeLivro(postlivro []byte) *PostarLivro {
	acao := PostarLivro{}
	posicao := 0
	acao.Epoca, posicao = util.ParseUint64(postlivro, posicao)
	acao.Autor, posicao = util.ParseToken(postlivro, posicao)
	// se o byte correspondente ao tipo de acao nao for o esperado, retorna nulo
	if postlivro[posicao] != APostarLivro {
		return nil
	}
	posicao += 1
	acao.TipoArquivo, posicao = util.ParseString(postlivro, posicao)
	acao.Conteudo, posicao = util.ParseHash(postlivro, posicao)
	acao.Data, posicao = util.ParseTime(postlivro, posicao)
	if posicao != len(postlivro) {
		return nil
	}
	return &acao
}
