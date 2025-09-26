// Não é parte do protocolo MEGA, apenas permite que a pessoa ingresse na rede
// do @s livres

package acoes

import (
	"github.com/freehandle/breeze/crypto"
	"github.com/freehandle/breeze/util"
)

type Entrar struct {
	Epoca   uint64
	Autor   crypto.Token
	Reasons string
}

func (c *Entrar) FazHash() crypto.Hash {
	return crypto.Hasher(c.Serializa())
}

func (c *Entrar) Affected() []crypto.Hash {
	return []crypto.Hash{crypto.ZeroHash}
}

func (c *Entrar) Autoria() crypto.Token {
	return c.Autor
}

func (c *Entrar) Serializa() []byte {
	bytes := make([]byte, 0)
	util.PutUint64(c.Epoca, &bytes)
	util.PutToken(c.Autor, &bytes)
	util.PutByte(AEntrarRede, &bytes)
	util.PutString(c.Reasons, &bytes)
	return bytes
}

func ParseSignIn(create []byte) *Entrar {
	action := Entrar{}
	position := 0
	action.Epoca, position = util.ParseUint64(create, position)
	action.Autor, position = util.ParseToken(create, position)
	if create[position] != AEntrarRede {
		return nil
	}
	position += 1
	action.Reasons, position = util.ParseString(create, position)
	if position > len(create) {
		return nil
	}
	return &action
}
