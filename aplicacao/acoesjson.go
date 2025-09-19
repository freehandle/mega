package aplicacao

import (
	"mega/protocolo/acoes"
	"time"
)

type PostarCauso struct {
	Acao       string    `json:"acao"`
	ID         int       `json:"id"`
	CampoCauso string    `json:"campoCauso"`
	DataHora   time.Time `json:"dataHora"`
}

func (a PostarCauso) ParaAcao() ([]acoes.Acao, error) {
	acao := acoes.PostarCauso{
		Conteudo: a.CampoCauso,
		Data:     a.DataHora,
	}
	return []acoes.Acao{&acao}, nil
}

type PostarFofoca struct {
	Acao        string    `json:"acao"`
	ID          int       `json:"id"`
	CampoFofoca string    `json:"campoFofoca"`
	DataHora    time.Time `json:"dataHora"`
}

func (a PostarFofoca) ParaAcao() ([]acoes.Acao, error) {
	acao := acoes.PostarFofoca{
		Conteudo: a.CampoFofoca,
		Data:     a.DataHora,
	}
	return []acoes.Acao{&acao}, nil
}

type PostarIdeia struct {
	Acao       string    `json:"acao"`
	ID         int       `json:"id"`
	CampoIdeia string    `json:"campoIdeia"`
	DataHora   time.Time `json:"dataHora"`
}

func (a PostarIdeia) ParaAcao() ([]acoes.Acao, error) {
	acao := acoes.PostarIdeia{
		Conteudo: a.CampoIdeia,
		Data:     a.DataHora,
	}
	return []acoes.Acao{&acao}, nil
}

type PostarLivro struct {
	Acao         string    `json:"acao"`
	ID           int       `json:"id"`
	TipoArquivo  string    `json:"tipoArquivo"`
	ArquivoLivro []byte    `json:"arquivoPraSubir"`
	DataHora     time.Time `json:"dataHora"`
}

func (a PostarLivro) ParaAcao() ([]acoes.Acao, error) {
	truncated := splitBytes(a.ArquivoLivro)
	allActions := make([]acoes.Acao, len(truncated.Parts))
	allActions[0] = &acoes.PostarLivro{
		Data:        a.DataHora,
		TipoArquivo: a.TipoArquivo,
	}
	for n := 1; n < len(truncated.Parts); n++ {
		allActions[n] = &acoes.MidiaMultiparte{
			Parte: byte(n),
			De:    byte(len(truncated.Parts)),
			Dados: truncated.Parts[n],
		}
	}
	return allActions, nil
}

type PostarMeme struct {
	Acao        string    `json:"acao"`
	ID          int       `json:"id"`
	TipoArquivo string    `json:"tipoArquivo"`
	ArquivoMeme []byte    `json:"arquivoPraSubir"`
	DataHora    time.Time `json:"dataHora"`
}

func (a PostarMeme) ParaAcao() ([]acoes.Acao, error) {
	truncated := splitBytes(a.ArquivoMeme)
	allActions := make([]acoes.Acao, len(truncated.Parts))
	allActions[0] = &acoes.PostarMeme{
		Data:        a.DataHora,
		TipoArquivo: a.TipoArquivo,
	}
	for n := 1; n < len(truncated.Parts); n++ {
		allActions[n] = &acoes.MidiaMultiparte{
			Parte: byte(n),
			De:    byte(len(truncated.Parts)),
			Dados: truncated.Parts[n],
		}
	}
	return allActions, nil
}

type PostarMusica struct {
	Acao        string    `json:"acao"`
	ID          int       `json:"id"`
	CampoMusica string    `json:"campoMusica"`
	DataHora    time.Time `json:"dataHora"`
}

func (a PostarMusica) ParaAcao() ([]acoes.Acao, error) {
	acao := acoes.PostarMusica{
		Conteudo: a.CampoMusica,
		Data:     a.DataHora,
	}
	return []acoes.Acao{&acao}, nil
}
