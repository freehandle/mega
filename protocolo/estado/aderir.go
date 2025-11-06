package estado

import (
	"time"

	"github.com/freehandle/breeze/crypto"
)

type Aderir struct {
	Autor crypto.Token
	Data  time.Time
	Hash  crypto.Hash
	// Apelido  string // vai ser fornecido pelo protocolo Apelidos
}

/*
// Verifica se a @ existe no @s livres
func (e *Estado) ValidaAderir(acao *acoes.AderirMIGA) error {
	if !e.VerificaSeMembro(acao.Autor) {
		return errors.New("autor do post nao reconhecido")
	}
	return

}
*/
