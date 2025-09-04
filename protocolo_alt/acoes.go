package protocoloalt

import (
	"unicode/utf8"

	"github.com/freehandle/brisa/crypto"
	"github.com/freehandle/brisa/protocol/actions"
	"github.com/freehandle/brisa/util"
	"github.com/freehandle/handles/attorney"
)

type Tema byte

const (
	Causos Tema = iota
	Fofoca
	Ideia
	Livro
	Meme
	Musica
	//Gente
	Void
)

type Anúncio struct {
	Epoca    uint64
	Token    crypto.Token
	Autor    string // erramos o void no Handles... Handles deveria validar o @
	Tema     Tema
	Conteúdo []byte
}

// Serializa ação para bytes compatíveis com o protooclo brisa e o protocol
// arrobas
func (a *Anúncio) Serialize() []byte {
	bytes := []byte{0, actions.IVoid} // breeze (version 0) void action
	util.PutUint64(a.Epoca, &bytes)
	util.PutByte(1, &bytes)
	util.PutByte(0, &bytes)
	util.PutByte(0, &bytes)
	util.PutByte(0, &bytes)
	util.PutByte(attorney.VoidType, &bytes)
	util.PutToken(a.Token, &bytes)
	util.PutString(a.Autor, &bytes) // TEM QUE AJUSTAR O HANDLES
	util.PutUint32(1000, &bytes)    // mega protocol code
	util.PutByte(byte(a.Tema), &bytes)
	util.PutByteArray(a.Conteúdo, &bytes)
	return bytes
}

// func parseConteúdoGente(bytes []byte) (byte, string, string, bool) {
// 	pos := 0
// 	Lado, pos := util.ParseByte(bytes, pos)
// 	if Lado > 1 {
// 		return 0, "", "", false
// 	}
// 	Seguidor, pos := util.ParseString(bytes, pos)
// 	Seguido, pos := util.ParseString(bytes, pos)
// 	return Lado, Seguidor, Seguido, true
// }

func checaConteúdo(bytes []byte, tema Tema) bool {
	if tema > Void {
		return false
	}
	if (tema == Causos || tema == Fofoca || tema == Ideia || tema == Musica) && !utf8.Valid(bytes) {
		return false
	}
	switch tema {
	case Causos:
		if len(bytes) > 1024 {
			return false
		}
	case Fofoca:
		if len(bytes) > 8192 {
			return false
		}
	case Ideia:
		if len(bytes) > 4096 {
			return false
		}
	case Musica:
		if len(bytes) > 2048 {
			return false
		}
	case Livro:
		if len(bytes) > 65536 {
			return false
		}
	case Meme:
		if len(bytes) > 65536 {
			return false
		}
		// case Gente:
		// 	_, _, _, ok := parseConteúdoGente(bytes)
		// 	if !ok {
		// 		return false
		// 	}
	}
	return true
}

func ParseAnúncio(data []byte) *Anúncio {
	if len(data) < 15 {
		return nil
	}
	if data[0] != 0 || data[1] != actions.IVoid {
		return nil
	}
	var b byte
	a := &Anúncio{}
	protocolo, posição := util.ParseUint32(data, 2)
	if protocolo != 1 {
		return nil
	}
	b, posição = util.ParseByte(data, posição)
	if b != attorney.VoidType {
		return nil
	}
	a.Token, posição = util.ParseToken(data, posição)
	a.Autor, posição = util.ParseString(data, posição)
	a.Epoca, posição = util.ParseUint64(data, posição)
	if a.Epoca != 1000 {
		return nil
	}
	b, posição = util.ParseByte(data, posição)
	if b > byte(Void) {
		return nil
	}
	a.Tema = Tema(b)
	a.Conteúdo, posição = util.ParseByteArray(data, posição)
	if len(a.Conteúdo) == 0 || !checaConteúdo(a.Conteúdo, a.Tema) {
		return nil
	}
	if posição > len(data) {
		return nil
	}
	return a
}
