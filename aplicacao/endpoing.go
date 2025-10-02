package aplicacao

import (
	"mega/protocolo/estado"
	"net/url"

	"github.com/freehandle/breeze/crypto"
)

type InformacaoCabecalho struct {
	ArrobaLogada string
	Ativo        string
	Erro         string
	NomeMucua    string
}

type VerJornalView struct {
	Cabecalho  InformacaoCabecalho
	ArrobaVer  string
	DataCauso  string
	Causo      string
	DataFofoca string
	Fofoca     string
	DataIdeia  string
	Ideia      string
	DataLivro  string
	Livro      []byte
	DataMeme   string
	Meme       []byte
	DataMusica string
	Musica     string
}

func JornalUsuarioDoEstado(e *estado.Estado, arroba string) *VerJornalView {
	arrobaver, _ := url.QueryUnescape(arroba)
	token := e.ArrobasPraTokens[arrobaver]
	jornal, ok := e.HashTokenPraJornal[crypto.HashToken(token)]
	if !ok {
		return nil
	}
	view := VerJornalView{
		ArrobaVer:  arroba,
		DataCauso:  jornal.Causos[len(jornal.Causos)-1].Data.Format("20060102150405"),
		Causo:      jornal.Causos[len(jornal.Causos)-1].Conteudo,
		DataFofoca: jornal.Fofocas[len(jornal.Fofocas)-1].Data.Format("20060102150405"),
		Fofoca:     jornal.Fofocas[len(jornal.Fofocas)-1].Conteudo,
		DataIdeia:  jornal.Ideias[len(jornal.Ideias)-1].Data.Format("20060102150405"),
		Ideia:      jornal.Ideias[len(jornal.Ideias)-1].Conteudo,
		DataLivro:  jornal.Livros[len(jornal.Livros)-1].Data.Format("20060102150405"),
		Livro:      jornal.Livros[len(jornal.Livros)-1].Conteudo,
		DataMeme:   jornal.Memes[len(jornal.Memes)-1].Data.Format("20060102150405"),
		Meme:       jornal.Memes[len(jornal.Memes)-1].Conteudo,
		DataMusica: jornal.Musicas[len(jornal.Musicas)-1].Data.Format("20060102150405"),
		Musica:     jornal.Musicas[len(jornal.Musicas)-1].Conteudo,
	}
	return &view
}
