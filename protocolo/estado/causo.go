package estado

import (
	"fmt"
	"time"
	"unicode/utf8"

	"github.com/freehandle/breeze/crypto"
)

type Causo struct {
	Conteudo string
	Autor    crypto.Token
	Data     time.Time
	Hash     crypto.Hash
}

// precisa ter ao menos 100 caracteres e no maximo dois paragrafos de caracteres (~1600 glyphos?)
func (i *Causo) ChecaFormato() bool {
	if utf8.RuneCountInString(i.Conteudo) < 100 || utf8.RuneCountInString(i.Conteudo) > 1600 {
		return false
	}
	return true
}

// o campo causo pode ser atualizado no maximo a cada 30 dias
func (i *Causo) ChecaTempo(s *Estado) bool {
	causos := s.HashTokenPraJornal[crypto.Hash(i.Autor)].Causos
	causoAntigo := causos[len(causos)-1]
	if causoAntigo != nil {
		if !i.Data.After(causoAntigo.Data) {
			fmt.Println("data antiga incorreta, verificar")
			return false

		}
		dias := i.Data.Sub(causoAntigo.Data).Hours() / 24
		if dias < 30 {
			fmt.Println("ainda não pode postar no campo causo")
			return false
		}
		// ja se passou um mes, pode postar
		return true
	}
	// nao tem post anterior de causo, pode postar
	return true
}
