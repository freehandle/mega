package indice

import (
	"github.com/freehandle/breeze/crypto"
	"github.com/freehandle/mega/protocolo/acoes"
)

type Jornal struct {
	Ideias  []*ConteudoData
	Memes   []*HashData
	Musicas []*ConteudoData
	Fofocas []*ConteudoData
	Causos  []*ConteudoData
	Livros  []*HashData
}

func NovoIndice() *Indice {
	return &Indice{
		ArrobaParaToken:  make(map[string]crypto.Token),
		TokenParaArroba:  make(map[crypto.Token]string),
		ArrobaParaJornal: make(map[string]*Jornal),
		HashToBytes:      make(map[crypto.Hash][]byte),
	}
}

type Indice struct {
	ArrobaParaToken  map[string]crypto.Token
	TokenParaArroba  map[crypto.Token]string
	ArrobaParaJornal map[string]*Jornal
	HashToBytes      map[crypto.Hash][]byte
}

type ConteudoData struct {
	Conteudo string
	Data     uint64
}

type HashData struct {
	Hash crypto.Hash
	Data uint64
}

func (i *Indice) IncorporaAutor(arroba string, token crypto.Token) {
	i.ArrobaParaToken[arroba] = token
	i.TokenParaArroba[token] = arroba
}

func (i *Indice) IncorporaConteudo(conteudo []byte) {
	hash := crypto.Hasher(conteudo)
	i.HashToBytes[hash] = conteudo
}

func (i *Indice) IncorporaCauso(causo *acoes.PostarCauso) {
	arroba := i.TokenParaArroba[causo.Autor]
	jornal, ok := i.ArrobaParaJornal[arroba]
	if !ok {
		jornal = &Jornal{}
	}
	novoCauso := &ConteudoData{
		Conteudo: causo.Conteudo,
		Data:     causo.Epoca,
	}
	if jornal.Causos == nil {
		jornal.Causos = []*ConteudoData{novoCauso}
	} else {
		jornal.Causos = append(jornal.Causos, novoCauso)
	}
}

func (i *Indice) IncorporaFofoca(causo *acoes.PostarFofoca) {
	arroba := i.TokenParaArroba[causo.Autor]
	jornal, ok := i.ArrobaParaJornal[arroba]
	if !ok {
		jornal = &Jornal{}
	}
	novoCauso := &ConteudoData{
		Conteudo: causo.Conteudo,
		Data:     causo.Epoca,
	}
	if jornal.Fofocas == nil {
		jornal.Fofocas = []*ConteudoData{novoCauso}
	} else {
		jornal.Fofocas = append(jornal.Fofocas, novoCauso)
	}
}

func (i *Indice) IncorporaIdeia(causo *acoes.PostarIdeia) {
	arroba := i.TokenParaArroba[causo.Autor]
	jornal, ok := i.ArrobaParaJornal[arroba]
	if !ok {
		jornal = &Jornal{}
	}
	novoCauso := &ConteudoData{
		Conteudo: causo.Conteudo,
		Data:     causo.Epoca,
	}
	if jornal.Ideias == nil {
		jornal.Ideias = []*ConteudoData{novoCauso}
	} else {
		jornal.Ideias = append(jornal.Ideias, novoCauso)
	}
}

func (i *Indice) IncorporaLivro(causo *acoes.PostarLivro) {
	arroba := i.TokenParaArroba[causo.Autor]
	jornal, ok := i.ArrobaParaJornal[arroba]
	if !ok {
		jornal = &Jornal{}
	}
	novoCauso := &HashData{
		Hash: causo.Conteudo,
		Data: causo.Epoca,
	}
	if jornal.Livros == nil {
		jornal.Livros = []*HashData{novoCauso}
	} else {
		jornal.Livros = append(jornal.Livros, novoCauso)
	}
}

func (i *Indice) IncorporaMeme(causo *acoes.PostarMeme) {
	arroba := i.TokenParaArroba[causo.Autor]
	jornal, ok := i.ArrobaParaJornal[arroba]
	if !ok {
		jornal = &Jornal{}
	}
	novoCauso := &HashData{
		Hash: causo.Conteudo,
		Data: causo.Epoca,
	}
	if jornal.Memes == nil {
		jornal.Memes = []*HashData{novoCauso}
	} else {
		jornal.Memes = append(jornal.Memes, novoCauso)
	}
}

func (i *Indice) IncorporaMusica(causo *acoes.PostarMusica) {
	arroba := i.TokenParaArroba[causo.Autor]
	jornal, ok := i.ArrobaParaJornal[arroba]
	if !ok {
		jornal = &Jornal{}
	}
	novoCauso := &ConteudoData{
		Conteudo: causo.Conteudo,
		Data:     causo.Epoca,
	}
	if jornal.Musicas == nil {
		jornal.Musicas = []*ConteudoData{novoCauso}
	} else {
		jornal.Musicas = append(jornal.Musicas, novoCauso)
	}
}

func (i *Indice) IncorporaAcao(dados []byte) {
	tipo := acoes.TipoDeAcao(dados)
	switch tipo {
	case acoes.APostarCauso:
		acao := acoes.LeCauso(dados)
		i.IncorporaCauso(acao)
	case acoes.APostarFofoca:
		acao := acoes.LeFofoca(dados)
		i.IncorporaFofoca(acao)
	case acoes.APostarIdeia:
		acao := acoes.LeIdeia(dados)
		// fmt.Printf("%+v\n", acao)
		i.IncorporaIdeia(acao)
	case acoes.APostarLivro:
		acao := acoes.LeLivro(dados)
		i.IncorporaLivro(acao)
	case acoes.APostarMeme:
		acao := acoes.LeMeme(dados)
		i.IncorporaMeme(acao)
	case acoes.APostarMusica:
		acao := acoes.LeMusica(dados)
		i.IncorporaMusica(acao)
	}
}
