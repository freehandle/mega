package protocoloalt

import (
	"github.com/freehandle/breeze/crypto"
	"github.com/freehandle/breeze/middleware/social"
)

type Proclame struct {
	Autor string
	Tema  Tema
}

type Elo struct {
	Seguidor string
	Seguido  string
}

type Mural struct {
	Proclames map[Proclame]uint64 // Proclame -> Época
	// Seguidos  map[string]uint64   // Seguido -> Época
	// Elos      map[Elo]struct{}    // Elo -> vazio
}

type Afixar struct {
	Proclames map[Proclame]uint64 // Proclame -> Época
	// Seguidos   map[string]uint64   // Seguido -> Época
	// NovosElos  map[Elo]struct{}    // Elo -> vazio
	// RemoveElos map[Elo]struct{}    // Elo -> vazio
}

func (a *Afixar) Merge(outras ...*Afixar) *Afixar {
	agrupadas := &Afixar{
		Proclames: make(map[Proclame]uint64),
		// Seguidos:   make(map[string]uint64),
		// NovosElos:  make(map[Elo]struct{}),
		// RemoveElos: make(map[Elo]struct{}),
	}
	return agrupadas
}

func (m *Mural) Validator(afixar ...*Afixar) *Afixando {
	if len(afixar) == 0 {
		return &Afixando{
			Mural:  m,
			Afixar: &Afixar{Proclames: make(map[Proclame]uint64)},
		}
	}
	if len(afixar) > 1 {
		afixar[0].Merge(afixar[1:]...)
	}
	return &Afixando{
		Mural:  m,
		Afixar: afixar[0],
	}
}

func (m *Mural) Incorporate(afixar *Afixar) {
	if afixar == nil {
		return
	}
	for proclame, época := range afixar.Proclames {
		m.Proclames[proclame] = época
	}
}

func (m *Mural) Checksum() crypto.Hash {
	// TODO: Checksum
	return crypto.ZeroHash
}

func (m *Mural) Recover() error {
	// TODO: Recovery
	return nil
}

func (m *Mural) Serialize() []byte {
	// TODO: Serialization
	return nil
}

func (m *Mural) Shutdown() {
}

func (m *Mural) Clone() chan social.Stateful[*Afixar, *Afixando] {
	cloned := make(chan social.Stateful[*Afixar, *Afixando], 2)
	clone := make(map[Proclame]uint64)
	for proclame, época := range m.Proclames {
		clone[proclame] = época
	}
	cloned <- &Mural{
		Proclames: clone,
	}
	return cloned
}
