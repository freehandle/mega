package estado

import (
	"errors"
	"time"

	"github.com/freehandle/breeze/crypto"
	"github.com/freehandle/mega/protocolo/acoes"
)

type Aderir struct {
	Autor crypto.Token
	Data  time.Time
	Hash  crypto.Hash
	// Apelido  string // vai ser fornecido pelo protocolo Apelidos
}

// Verifica se a @ existe no @s livres
func (e *Estado) ValidaAderir(acao *acoes.AderirMIGA) error {
	if !e.VerificaSeMembro(acao.Autor) {
		return errors.New("autor do post nao reconhecido")
	}
	return

}
