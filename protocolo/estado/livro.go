package estado

import (
	"fmt"
	"time"

	"github.com/freehandle/breeze/crypto"
)

type Livro struct {
	Conteudo []byte
	Autor    crypto.Token //arroba
	Data     time.Time
}

// precisa ser uma imagem
func (i *Livro) ChecaFormato() bool {
	// ver aqui o que eu faco
	return true
}

// o campo livro pode ser atualizado no maximo a cada 30 dias
func (i *Livro) ChecaTempo(s *Estado) bool {
	livroAntigo := s.HashTokenPraJornal[crypto.Hash(i.Autor)].livro
	if livroAntigo != nil {
		if !i.Data.After(livroAntigo.Data) {
			fmt.Println("data antiga incorreta, verificar")
			return false

		}
		dias := i.Data.Sub(livroAntigo.Data).Hours() / 24
		if dias < 7 {
			fmt.Println("ainda nao pode postar")
			return false
		}
		// ja se passou um mes, pode postar
		return true
	}
	// nao tem post anterior de livro, pode postar
	return true
}
