package estado

import (
	"sync"

	"github.com/freehandle/breeze/crypto"
	"github.com/freehandle/breeze/middleware/social"
	"github.com/freehandle/mega/protocolo/acoes"
)

const (
	LapsoCauso  = 30 * 24 * 60 * 60
	LapsoFofoca = 15 * 24 * 60 * 60
	LapsoIdeia  = 30 * 24 * 60 * 60
	LapsoLivro  = 30 * 24 * 60 * 60
	LapsoMeme   = 7 * 24 * 60 * 60
	LapsoMusica = 15 * 24 * 60 * 60
)

type Postagem struct {
	Token crypto.Token
	Tipo  byte
}

func Genesis(epoch uint64) *Estado {
	return &Estado{
		mu:                sync.Mutex{},
		Epoca:             epoch,
		UltimaAtualizacao: make(map[Postagem]uint64),
	}
}

type Estado struct {
	mu                sync.Mutex
	Epoca             uint64
	UltimaAtualizacao map[Postagem]uint64
}

type EstadoMutante struct {
	Original *Estado
	Mutacoes *Mutacoes
}

func (e *EstadoMutante) Mutations() *Mutacoes {
	return e.Mutacoes
}

func (e *EstadoMutante) UltimaPublicao(token crypto.Token, tipo byte) uint64 {
	postagem := Postagem{Token: token, Tipo: tipo}
	if ultima, ok := e.Mutacoes.Atualizacoes[postagem]; ok {
		return ultima
	}
	if ultima, ok := e.Original.UltimaAtualizacao[postagem]; ok {
		return ultima
	}
	return 0
}

func (e *EstadoMutante) Validate(dados []byte) bool {
	tipo := acoes.TipoDeAcao(dados)
	switch tipo {
	case acoes.APostarCauso:
		acao := acoes.LeCauso(dados)
		if acao == nil || !acao.ValidarFormato() {
			return false
		}
		ultima := e.UltimaPublicao(acao.Autor, acoes.APostarCauso)
		if ultima > 0 && (acao.Epoca < ultima+LapsoCauso) {
			return false
		}
		if acao.Epoca <= e.Mutacoes.Epoca {
			postagem := Postagem{Token: acao.Autor, Tipo: acoes.APostarCauso}
			e.Mutacoes.Atualizacoes[postagem] = acao.Epoca
			return true
		}
		return false
	case acoes.APostarFofoca:
		acao := acoes.LeFofoca(dados)
		if acao == nil || !acao.ValidarFormato() {
			return false
		}
		ultima := e.UltimaPublicao(acao.Autor, acoes.APostarFofoca)
		if ultima > 0 && (acao.Epoca < ultima+LapsoFofoca) {
			return false
		}
		if acao.Epoca <= e.Mutacoes.Epoca {
			postagem := Postagem{Token: acao.Autor, Tipo: acoes.APostarFofoca}
			e.Mutacoes.Atualizacoes[postagem] = acao.Epoca
			return true
		}
		return false
	case acoes.APostarIdeia:
		acao := acoes.LeIdeia(dados)
		if acao == nil || !acao.ValidarFormato() {
			return false
		}
		ultima := e.UltimaPublicao(acao.Autor, acoes.APostarIdeia)
		if ultima > 0 && (acao.Epoca < ultima+LapsoIdeia) {
			return false
		}
		if acao.Epoca <= e.Mutacoes.Epoca {
			postagem := Postagem{Token: acao.Autor, Tipo: acoes.APostarIdeia}
			e.Mutacoes.Atualizacoes[postagem] = acao.Epoca
			return true
		}
		return false
	case acoes.APostarLivro:
		acao := acoes.LeLivro(dados)
		if acao == nil || !acao.ValidarFormato() {
			return false
		}
		ultima := e.UltimaPublicao(acao.Autor, acoes.APostarLivro)
		if ultima > 0 && (acao.Epoca < ultima+LapsoLivro) {
			return false
		}
		if acao.Epoca <= e.Mutacoes.Epoca {
			postagem := Postagem{Token: acao.Autor, Tipo: acoes.APostarLivro}
			e.Mutacoes.Atualizacoes[postagem] = acao.Epoca
			return true
		}
		return false
	case acoes.APostarMeme:
		acao := acoes.LeMeme(dados)
		if acao == nil || !acao.ValidarFormato() {
			return false
		}
		ultima := e.UltimaPublicao(acao.Autor, acoes.APostarMeme)
		if ultima > 0 && (acao.Epoca < ultima+LapsoMeme) {
			return false
		}
		if acao.Epoca <= e.Mutacoes.Epoca {
			postagem := Postagem{Token: acao.Autor, Tipo: acoes.APostarMeme}
			e.Mutacoes.Atualizacoes[postagem] = acao.Epoca
			return true
		}
		return false
	case acoes.APostarMusica:
		acao := acoes.LeMusica(dados)
		if acao == nil || !acao.ValidarFormato() {
			return false
		}
		ultima := e.UltimaPublicao(acao.Autor, acoes.APostarMusica)
		if ultima > 0 && (acao.Epoca < ultima+LapsoMusica) {
			return false
		}
		if acao.Epoca <= e.Mutacoes.Epoca {
			postagem := Postagem{Token: acao.Autor, Tipo: acoes.APostarMusica}
			e.Mutacoes.Atualizacoes[postagem] = acao.Epoca
			return true
		}
		return false
	default:
		return false
	}
}

func (e *Estado) Validator(mutacoes ...*Mutacoes) *EstadoMutante {
	validador := EstadoMutante{
		Original: e,
	}
	if len(mutacoes) == 0 {
		validador.Mutacoes = &Mutacoes{
			Epoca:        e.Epoca + 1,
			Atualizacoes: make(map[Postagem]uint64),
		}
	} else if len(mutacoes) == 1 {
		validador.Mutacoes = mutacoes[0]
	} else {
		validador.Mutacoes = mutacoes[0].Merge(mutacoes[1:]...)
	}
	return &validador
}

func (e *Estado) Incorporate(mutacoes *Mutacoes) {
	e.Epoca = mutacoes.Epoca
	for k, v := range mutacoes.Atualizacoes {
		e.UltimaAtualizacao[k] = v
	}
}

func (e *Estado) Shutdown() {

}

func (e *Estado) Checksum() crypto.Hash {
	return crypto.ZeroHash
}

func (e *Estado) Clone() chan social.Stateful[*Mutacoes, *EstadoMutante] {
	e.mu.Lock()
	defer e.mu.Unlock()
	clone := Estado{
		Epoca:             e.Epoca,
		UltimaAtualizacao: make(map[Postagem]uint64),
	}
	for k, v := range e.UltimaAtualizacao {
		clone.UltimaAtualizacao[k] = v
	}
	resposta := make(chan social.Stateful[*Mutacoes, *EstadoMutante], 2)
	resposta <- &clone
	return resposta
}

func (e *Estado) Serialize() []byte {
	return nil
}

type Mutacoes struct {
	Epoca        uint64
	Atualizacoes map[Postagem]uint64
}

func (m *Mutacoes) Merge(mutacoes ...*Mutacoes) *Mutacoes {
	merged := &Mutacoes{
		Atualizacoes: make(map[Postagem]uint64),
	}
	for _, mu := range append([]*Mutacoes{m}, mutacoes...) {
		for k, v := range mu.Atualizacoes {
			if at, ok := merged.Atualizacoes[k]; !ok || v > at {
				merged.Atualizacoes[k] = v
			}
		}
	}
	return merged
}

/*var TiposImagens = []string{"jpg", "gif", "png", "bmp", "svg"}

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
*/
