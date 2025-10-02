package estado

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/freehandle/breeze/crypto"
)

type Livro struct {
	Conteudo []byte
	Autor    crypto.Token
	Data     time.Time
	Hash     crypto.Hash
}

// precisa ser uma imagem
func (i *Livro) ChecaFormato(tipo string) bool {
	tipomin := strings.ToLower(tipo)
	return slices.Contains(TiposImagens, tipomin)
}

// o campo livro pode ser atualizado no maximo a cada 30 dias
func (i *Livro) ChecaTempo(s *Estado) bool {
	livros := s.HashTokenPraJornal[crypto.Hash(i.Autor)].Livros
	livroAntigo := livros[len(livros)-1]
	if livroAntigo != nil {
		if !i.Data.After(livroAntigo.Data) {
			fmt.Println("data antiga incorreta, verificar")
			return false

		}
		dias := i.Data.Sub(livroAntigo.Data).Hours() / 24
		if dias < 30 {
			fmt.Println("ainda não pode postar no campo livro")
			return false
		}
		// ja se passou um mes, pode postar
		return true
	}
	// nao tem post anterior de livro, pode postar
	return true
}
