package estado

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/freehandle/breeze/crypto"
)

type Meme struct {
	Conteudo []byte
	Autor    crypto.Token
	Data     time.Time
	Hash     crypto.Hash
}

// precisa ser uma imagem
func (i *Meme) ChecaFormato(tipo string) bool {
	tipomin := strings.ToLower(tipo)
	return slices.Contains(TiposImagens, tipomin)
}

// o campo meme pode ser atualizado no maximo a cada 7 dias
func (i *Meme) ChecaTempo(s *Estado) bool {
	memes := s.HashTokenPraJornal[crypto.Hash(i.Autor)].Memes
	memeAntigo := memes[len(memes)-1]
	if memeAntigo != nil {
		if !i.Data.After(memeAntigo.Data) {
			fmt.Println("data antiga incorreta, verificar")
			return false

		}
		dias := i.Data.Sub(memeAntigo.Data).Hours() / 24
		if dias < 7 {
			fmt.Println("ainda não pode postar no campo meme")
			return false
		}
		// ja se passou um mes, pode postar
		return true
	}
	// nao tem post anterior de meme, pode postar
	return true
}
