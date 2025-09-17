package estado

import (
	"fmt"
	"time"
	"unicode/utf8"

	"github.com/freehandle/breeze/crypto"
)

type Fofoca struct {
	Conteudo string
	Autor    crypto.Token
	Data     time.Time
	Hash     crypto.Hash
}

// precisa ter ao menos 100 caracteres e no maximo duas paginas de caracteres (~2500 glyphos?)
func (i *Fofoca) ChecaFormato() bool {
	if utf8.RuneCountInString(i.Conteudo) < 100 || utf8.RuneCountInString(i.Conteudo) > 5000 {
		return false
	}
	return true
}

// o campo fofoca pode ser atualizado no maximo a cada 15 dias
func (i *Fofoca) ChecaTempo(s *Estado) bool {
	fofocas := s.HashTokenPraJornal[crypto.Hash(i.Autor)].fofoca
	fofocaAntiga := fofocas[len(fofocas)-1]
	if fofocaAntiga != nil {
		if !i.Data.After(fofocaAntiga.Data) {
			fmt.Println("data antiga incorreta, verificar")
			return false

		}
		dias := i.Data.Sub(fofocaAntiga.Data).Hours() / 24
		if dias < 15 {
			fmt.Println("ainda não pode postar no campo fofoca")
			return false
		}
		// ja se passou uma quinzena, pode postar
		return true
	}
	// nao tem post anterior de fofoca, pode postar
	return true
}
