package estado

import (
	"fmt"
	"time"
	"unicode/utf8"

	"github.com/freehandle/breeze/crypto"
)

type Musica struct {
	Conteudo string
	Autor    crypto.Token //arroba
	Data     time.Time
}

// precisa ter ao menos 10 caracteres e no maximo um paragrafo (~800 glyphos?)
func (i *Musica) ChecaFormato() bool {
	if utf8.RuneCountInString(i.Conteudo) < 10 || utf8.RuneCountInString(i.Conteudo) > 800 {
		return false
	}
	return true
}

// o campo musica pode ser atualizado no maximo a cada 15 dias
func (i *Musica) ChecaTempo(s *Estado) bool {
	musicaAntiga := s.HashTokenPraJornal[crypto.Hash(i.Autor)].musica
	if musicaAntiga != nil {
		if !i.Data.After(musicaAntiga.Data) {
			fmt.Println("data antiga incorreta, verificar")
			return false

		}
		dias := i.Data.Sub(musicaAntiga.Data).Hours() / 24
		if dias < 15 {
			fmt.Println("ainda nao pode postar")
			return false
		}
		// ja se passou uma quinzena, pode postar
		return true
	}
	// nao tem post anterior de musica, pode postar
	return true
}
