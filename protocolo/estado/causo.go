package estado

import (
	"fmt"
	"time"
	"unicode/utf8"

	"github.com/freehandle/breeze/crypto"
)

type Causo struct {
	Conteudo string
	Autor    crypto.Token //arroba
	Data     time.Time
}

// precisa ter ao menos 1 caractere e no maximo duas laudas de caracteres (2500/pag, approx 2pg?)
func (i *Causo) ChecaFormato() bool {
	if utf8.RuneCountInString(i.Conteudo) < 100 || utf8.RuneCountInString(i.Conteudo) > 5000 {
		return false
	}
	return true
}

// o campo causo pode ser atualizado no maximo a cada 30 dias
func (i *Causo) ChecaTempo(s *Estado) bool {
	causoAntigo := s.HashTokenPraJornal[crypto.Hash(i.Autor)].causo
	if causoAntigo != nil {
		if !i.Data.After(causoAntigo.Data) {
			fmt.Println("data antiga incorreta, verificar")
			return false

		}
		dias := i.Data.Sub(causoAntigo.Data).Hours() / 24
		if dias < 30 {
			fmt.Println("ainda nao pode postar")
			return false
		}
		// ja se passou um mes, pode postar
		return true
	}
	// nao tem post anterior de causo, pode postar
	return true
}
