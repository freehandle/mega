package acoes

import (
	"github.com/freehandle/breeze/crypto"
	"github.com/freehandle/breeze/util"
)

// Implementa a quebra de um arquivo longo em multiplas acoes com cada parte do arquivo pra ele subir na blockchain
type MidiaMultiparte struct {
	Epoca uint64
	Autor crypto.Token
	Parte byte
	De    byte
	Dados []byte
}

func (m *MidiaMultiparte) FazHash() crypto.Hash {
	return crypto.Hasher([]byte(m.Serializa()))
}

func (m *MidiaMultiparte) Autoria() crypto.Token {
	return m.Autor
}

func (m *MidiaMultiparte) Serializa() []byte {
	bytes := make([]byte, 0)
	util.PutUint64(m.Epoca, &bytes)
	util.PutToken(m.Autor, &bytes)
	util.PutByte(AMidiaMultiparte, &bytes)
	util.PutByte(m.Parte, &bytes)
	util.PutByte(m.De, &bytes)
	util.PutByteArray(m.Dados, &bytes)
	return bytes
}

func LeMidiaMultiparte(create []byte) *MidiaMultiparte {
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
