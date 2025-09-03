package estado

import "github.com/freehandle/breeze/crypto"

type ArrobaUsuario struct {
	arroba string
}

type Estado struct {
	Epoca               uint64
	ArrobasPraTokens    map[string]crypto.Token
	HashTokenPraArrobas map[crypto.Hash]string
	HashTokenPraJornal  map[crypto.Hash]*Jornal
}

type Jornal struct {
	ideia  *Ideia
	meme   *Meme
	musica *Musica
	fofoca *Fofoca
	causo  *Causo
	livro  *Livro
}
