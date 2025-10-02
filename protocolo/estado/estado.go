package estado

import (
	"errors"
	"time"

	"github.com/freehandle/mega/protocolo/acoes"

	"github.com/freehandle/breeze/crypto"
)

var TiposImagens = []string{"jpg", "gif", "pgn", "bmp", "svg"}

type InfoUsuario struct {
	Arroba string
}

type EncontraApelido interface {
	Arroba(token crypto.Token) *InfoUsuario
	Token(arroba string) *crypto.Token
}

type Estado struct {
	Epoca               uint64
	ArrobasPraTokens    map[string]crypto.Token
	HashTokenPraArrobas map[crypto.Hash]string
	HashTokenPraJornal  map[crypto.Hash]*Jornal
	Apelidos            EncontraApelido
	HoraGenesis         time.Time
}

type Jornal struct {
	Ideias  []*Ideia
	Memes   []*Meme
	Musicas []*Musica
	Fofocas []*Fofoca
	Causos  []*Causo
	Livros  []*Livro
}

// inicializa o estado do protocolo mega
func EstadoInicial() *Estado {
	estado := &Estado{
		Epoca:               0,
		ArrobasPraTokens:    make(map[string]crypto.Token),
		HashTokenPraArrobas: make(map[crypto.Hash]string),
		HashTokenPraJornal:  map[crypto.Hash]*Jornal{},
	}
	return estado
}

// verifica se um token fornecido é membro da rede
func (e *Estado) VerificaSeMembro(token crypto.Token) bool {
	hash := crypto.HashToken(token)
	_, katu := e.HashTokenPraArrobas[hash]
	return katu // katu quer dizer bom, certo, ok, em tupinambá :)
}

// Le e valida as ações do protocolo
func (e *Estado) Acao(dados []byte) error {
	tipo := acoes.TipoDeAcao(dados)
	switch tipo {
	case acoes.APostarCauso:
		acao := acoes.LeCauso(dados)
		if acao == nil {
			return errors.New("nao foi possivel ler a acao do tipo post causo")
		}
		err := e.ValidaCauso(acao)
		return err
	case acoes.APostarFofoca:
		acao := acoes.LeFofoca(dados)
		if acao == nil {
			return errors.New("nao foi possivel ler a acao do tipo post fofoca")
		}
		err := e.ValidaFofoca(acao)
		return err
	case acoes.APostarIdeia:
		acao := acoes.LeIdeia(dados)
		if acao == nil {
			return errors.New("nao foi possivel ler a acao do tipo post ideia")
		}
		err := e.ValidaIdeia(acao)
		return err
	case acoes.APostarLivro:
		acao := acoes.LeLivro(dados)
		if acao == nil {
			return errors.New("nao foi possivel ler a acao do tipo post livro")
		}
		err := e.ValidaLivro(acao)
		return err
	case acoes.APostarMeme:
		acao := acoes.LeMeme(dados)
		if acao == nil {
			return errors.New("nao foi possivel ler a acao do tipo post meme")
		}
		err := e.ValidaMeme(acao)
		return err
	case acoes.APostarMusica:
		acao := acoes.LeMusica(dados)
		if acao == nil {
			return errors.New("nao foi possivel ler a acao do tipo post musica")
		}
		err := e.ValidaMusica(acao)
		return err
	}
	return errors.New("acao nao reconhecida")
}

// Valida ação postar causo
func (e *Estado) ValidaCauso(acao *acoes.PostarCauso) error {
	if !e.VerificaSeMembro(acao.Autor) {
		return errors.New("autor do post nao reconhecido")
	}
	novocauso := Causo{
		Conteudo: acao.Conteudo,
		Autor:    acao.Autor,
		Data:     acao.Data,
		Hash:     acao.FazHash(),
	}
	if !novocauso.ChecaFormato() {
		return errors.New("formato do causo nao esta adequado")
	}
	if !novocauso.ChecaTempo(e) {
		return errors.New("ainda nao pode postar no campo causo")
	}
	e.HashTokenPraJornal[crypto.Hash(acao.Autor)].Causos = append(e.HashTokenPraJornal[crypto.Hash(acao.Autor)].Causos, &novocauso)
	return nil
}

// Valida a ação postar fofoca
func (e *Estado) ValidaFofoca(acao *acoes.PostarFofoca) error {
	if !e.VerificaSeMembro(acao.Autor) {
		return errors.New("autor do post nao reconhecido")
	}
	novafofoca := Fofoca{
		Conteudo: acao.Conteudo,
		Autor:    acao.Autor,
		Data:     acao.Data,
		Hash:     acao.FazHash(),
	}
	if !novafofoca.ChecaFormato() {
		return errors.New("formato da fofoca nao esta adequado")
	}
	if !novafofoca.ChecaTempo(e) {
		return errors.New("ainda nao pode postar no campo fofoca")
	}
	e.HashTokenPraJornal[crypto.Hash(acao.Autor)].Fofocas = append(e.HashTokenPraJornal[crypto.Hash(acao.Autor)].Fofocas, &novafofoca)
	return nil
}

// Valida a ação postar ideia
func (e *Estado) ValidaIdeia(acao *acoes.PostarIdeia) error {
	if !e.VerificaSeMembro(acao.Autor) {
		return errors.New("autor do post nao reconhecido")
	}
	novaideia := Ideia{
		Conteudo: acao.Conteudo,
		Autor:    acao.Autor,
		Data:     acao.Data,
		Hash:     acao.FazHash(),
	}
	if !novaideia.ChecaFormato() {
		return errors.New("formado da ideia nao esta adequado")
	}
	if !novaideia.ChecaTempo(e) {
		return errors.New("ainda nao pode postar no campo ideia")
	}
	e.HashTokenPraJornal[crypto.Hash(acao.Autor)].Ideias = append(e.HashTokenPraJornal[crypto.Hash(acao.Autor)].Ideias, &novaideia)
	return nil
}

// Valida a ação postar livro
func (e *Estado) ValidaLivro(acao *acoes.PostarLivro) error {
	if !e.VerificaSeMembro(acao.Autor) {
		return errors.New("autor do post nao reconhecido")
	}
	novolivro := Livro{
		Conteudo: acao.Conteudo,
		Autor:    acao.Autor,
		Data:     acao.Data,
		Hash:     acao.FazHash(),
	}
	if !novolivro.ChecaFormato(acao.TipoArquivo) {
		return errors.New("formado do livro nao esta adequado")
	}
	if !novolivro.ChecaTempo(e) {
		return errors.New("ainda nao pode postar no campo livro")
	}
	e.HashTokenPraJornal[crypto.Hash(acao.Autor)].Livros = append(e.HashTokenPraJornal[crypto.Hash(acao.Autor)].Livros, &novolivro)
	return nil
}

// Valida a ção postar meme
func (e *Estado) ValidaMeme(acao *acoes.PostarMeme) error {
	if !e.VerificaSeMembro(acao.Autor) {
		return errors.New("autor do post nao reconhecido")
	}
	novomeme := Meme{
		Conteudo: acao.Conteudo,
		Autor:    acao.Autor,
		Data:     acao.Data,
		Hash:     acao.FazHash(),
	}
	if !novomeme.ChecaFormato(acao.TipoArquivo) {
		return errors.New("formato do meme nao esta adequado")
	}
	if !novomeme.ChecaTempo(e) {
		return errors.New("ainda nao pode postar no campo meme")
	}
	e.HashTokenPraJornal[crypto.Hash(acao.Autor)].Memes = append(e.HashTokenPraJornal[crypto.Hash(acao.Autor)].Memes, &novomeme)
	return nil
}

// Valida a ação postar música
func (e *Estado) ValidaMusica(acao *acoes.PostarMusica) error {
	if !e.VerificaSeMembro(acao.Autor) {
		return errors.New("autor do post nao reconhecido")
	}
	novamusica := Musica{
		Conteudo: acao.Conteudo,
		Autor:    acao.Autor,
		Data:     acao.Data,
		Hash:     acao.FazHash(),
	}
	if !novamusica.ChecaFormato() {
		return errors.New("formato da musica nao esta adequado")
	}
	if !novamusica.ChecaTempo(e) {
		return errors.New("ainda nao pode postar no campo musica")
	}
	e.HashTokenPraJornal[crypto.Hash(acao.Autor)].Musicas = append(e.HashTokenPraJornal[crypto.Hash(acao.Autor)].Musicas, &novamusica)
	return nil
}
