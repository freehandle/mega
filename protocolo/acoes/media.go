package acoes

import (
	"github.com/freehandle/breeze/crypto"
	"github.com/freehandle/breeze/util"
)

type MidiaMultiparte struct {
	Epoca uint64
	Autor crypto.Token
	Parte byte
	De    byte
	Dados []byte
}

func (c *MidiaMultiparte) Authored() crypto.Token {
	return c.Autor
}

func (c *MidiaMultiparte) Serialize() []byte {
	bytes := make([]byte, 0)
	util.PutUint64(c.Epoca, &bytes)
	util.PutToken(c.Autor, &bytes)
	util.PutByte(AMidiaMultiparte, &bytes)
	util.PutByte(c.Parte, &bytes)
	util.PutByte(c.De, &bytes)
	util.PutByteArray(c.Dados, &bytes)
	return bytes
}

func ParseMultipartMedia(create []byte) *MidiaMultiparte {
	action := MidiaMultiparte{}
	position := 0
	action.Epoca, position = util.ParseUint64(create, position)
	action.Autor, position = util.ParseToken(create, position)
	if create[position] != AMidiaMultiparte {
		return nil
	}
	position += 1
	action.Parte, position = util.ParseByte(create, position)
	action.De, position = util.ParseByte(create, position)
	action.Dados, position = util.ParseByteArray(create, position)
	if position != len(create) {
		return nil
	}
	return &action
}
