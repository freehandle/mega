package estado

import (
	"fmt"
	"time"
	"unicode/utf8"

	"github.com/freehandle/breeze/crypto"
)

type Ideia struct {
	Conteudo string
	Autor    crypto.Token
	Data     time.Time
	Hash     crypto.Hash
}

// precisa ter ao menos 1 caractere e no maximo uma pagina de caracteres (~2500 glyphos?)
func (i *Ideia) ChecaFormato() bool {
	if utf8.RuneCountInString(i.Conteudo) < 1 || utf8.RuneCountInString(i.Conteudo) > 2500 {
		return false
	}
	return true
}

// o campo ideia pode ser atualizado no maximo a cada 30 dias
func (i *Ideia) ChecaTempo(s *Estado) bool {
	ideias := s.HashTokenPraJornal[crypto.Hash(i.Autor)].ideia
	ideiaAntiga := ideias[len(ideias)-1]
	if ideiaAntiga != nil {
		if !i.Data.After(ideiaAntiga.Data) {
			fmt.Println("data antiga incorreta, verificar")
			return false

		}
		dias := i.Data.Sub(ideiaAntiga.Data).Hours() / 24
		if dias < 30 {
			fmt.Println("ainda não pode postar no campo ideia")
			return false
		}
		// ja se passou um mes, pode postar
		return true
	}
	// nao tem post anterior de ideia, pode postar
	return true
}
