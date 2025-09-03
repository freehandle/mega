package estado

import (
	"fmt"
	"time"

	"github.com/freehandle/breeze/crypto"
)

type Meme struct {
	Conteudo []byte
	Autor    crypto.Token //arroba
	Data     time.Time
}

// precisa ser uma imagem
func (i *Meme) ChecaFormato() bool {
	// ver aqui o que eu faco
	return true
}

// o campo meme pode ser atualizado no maximo a cada 30 dias
func (i *Meme) ChecaTempo(s *Estado) bool {
	memeAntigo := s.HashTokenPraJornal[crypto.Hash(i.Autor)].meme
	if memeAntigo != nil {
		if !i.Data.After(memeAntigo.Data) {
			fmt.Println("data antiga incorreta, verificar")
			return false

		}
		dias := i.Data.Sub(memeAntigo.Data).Hours() / 24
		if dias < 7 {
			fmt.Println("ainda nao pode postar")
			return false
		}
		// ja se passou um mes, pode postar
		return true
	}
	// nao tem post anterior de meme, pode postar
	return true
}
