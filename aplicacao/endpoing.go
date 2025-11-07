package aplicacao

/*
import (
	"net/url"
	"time"

	"github.com/freehandle/mega/protocolo/estado"

	"github.com/freehandle/breeze/crypto"
)

type InformacaoCabecalho struct {
	ArrobaLog       string
	Ativo           string
	Erro            string
	NomeMucua       string
	LinkSelecionada string
}

type ViewConvite struct {
	Cabecalho InformacaoCabecalho
	Seed      string
	// pra testar
	Nome  string
	Nome2 string
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
		// return nil
		view := VerJornalView{
			ArrobaVer:  "teste",
			DataCauso:  time.Time.Format(time.Now(), "2006-01-02 às 15h04m"),
			Causo:      "TESTE",
			DataFofoca: time.Time.Format(time.Now(), "2006-01-02 às 15h04m"),
			Fofoca:     "NOVO TESTE",
			DataIdeia:  time.Time.Format(time.Now(), "2006-01-02 às 15h04m"),
			Ideia:      "NOVO TESTE",
			DataLivro:  time.Time.Format(time.Now(), "2006-01-02 às 15h04m"),
			Livro:      make([]byte, 2),
			DataMeme:   time.Time.Format(time.Now(), "2006-01-02 às 15h04m"),
			Meme:       make([]byte, 2),
			DataMusica: time.Time.Format(time.Now(), "2006-01-02 às 15h04m"),
			Musica:     "AAAAAAAAAAAAAAAAAAAAAA",
		}
		return &view
	}
	view := VerJornalView{
		ArrobaVer:  arroba,
		DataCauso:  jornal.Causos[len(jornal.Causos)-1].Data.Format("2006-01-02 às 15h04m"),
		Causo:      jornal.Causos[len(jornal.Causos)-1].Conteudo,
		DataFofoca: jornal.Fofocas[len(jornal.Fofocas)-1].Data.Format("2006-01-02 às 15h04m"),
		Fofoca:     jornal.Fofocas[len(jornal.Fofocas)-1].Conteudo,
		DataIdeia:  jornal.Ideias[len(jornal.Ideias)-1].Data.Format("2006-01-02 às 15h04m"),
		Ideia:      jornal.Ideias[len(jornal.Ideias)-1].Conteudo,
		DataLivro:  jornal.Livros[len(jornal.Livros)-1].Data.Format("2006-01-02 às 15h04m"),
		Livro:      jornal.Livros[len(jornal.Livros)-1].Conteudo,
		DataMeme:   jornal.Memes[len(jornal.Memes)-1].Data.Format("2006-01-02 às 15h04m"),
		Meme:       jornal.Memes[len(jornal.Memes)-1].Conteudo,
		DataMusica: jornal.Musicas[len(jornal.Musicas)-1].Data.Format("2006-01-02 às 15h04m"),
		Musica:     jornal.Musicas[len(jornal.Musicas)-1].Conteudo,
	}
	return &view
}

func AderirProtocoloMIGA(e *estado.Estado, token crypto.Token, arroba string) {
	// FAZER VERIFICACOES USANDO IU
	e.ArrobasPraTokens[arroba] = token
	e.HashTokenPraArrobas[crypto.HashToken(token)] = arroba
	e.HashTokenPraJornal[crypto.HashToken(token)] = Jornal{}
}

func PostarTexto(e *estado.Estado, arroba string, tipoAcao string) {
	arrobaver, _ := url.QueryUnescape(arroba)
	token := e.ArrobasPraTokens[arrobaver]
	return
	// jornal, ok := e.HashTokenPraJornal[crypto.HashToken(token)]
	// if !ok {
	// 	e.HashTokenPraJornal
	// }
}
*/
