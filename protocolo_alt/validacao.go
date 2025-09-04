package protocoloalt

type Afixando struct {
	Mural  *Mural
	Afixar *Afixar
}

func (a *Afixando) Mutations() *Afixar {
	return a.Afixar
}

func (a *Afixando) Validate(data []byte) bool {
	anúncio := ParseAnúncio(data)
	if anúncio == nil || a.Mural == nil || a.Afixar == nil {
		return false
	}
	proclame := Proclame{Autor: anúncio.Autor, Tema: anúncio.Tema}
	mural := a.Mural.Proclames[proclame]
	afixando := a.Afixar.Proclames[proclame]
	prazo := uint64(0)
	switch anúncio.Tema {
	case Causos:
		prazo = PrazoCausos
	case Fofoca:
		prazo = PrazoFofoca
	case Ideia:
		prazo = PrazoIdeia
	case Livro:
		prazo = PrazoLivro
	case Meme:
		prazo = PrazoMeme
	case Musica:
		prazo = PrazoMusica
		// case Gente:
		// 	prazo = PrazoGente
	}
	if (anúncio.Epoca-prazo < mural) || (anúncio.Epoca-prazo < afixando) {
		return false
	}
	a.Afixar.Proclames[proclame] = anúncio.Epoca
	return true
}
