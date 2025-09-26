package aplicacao

import (
	"errors"
	"mega/protocolo/estado"
	"slices"
)

type Index struct {
	Membros    []string
	Seguindo   map[string][]string
	Seguidores map[string][]string
	estado     *estado.Estado
}

func NovoIndex() *Index {
	return &Index{
		Membros:    make([]string, 0),
		Seguindo:   make(map[string][]string),
		Seguidores: make(map[string][]string),
	}
}

func (i *Index) ChecaSeMembro(arroba string) bool {
	return slices.Contains(i.Membros, arroba)
}

func (i *Index) Seguir(autor string, alvo string) error {
	if !i.ChecaSeMembro(autor) {
		return errors.New("autor nao é membro")
	}
	if !i.ChecaSeMembro(alvo) {
		return errors.New("alvo nao é membro")
	}
	if listaAutor, fazparteAutor := i.Seguindo[autor]; fazparteAutor {
		if !slices.Contains(i.Seguindo[autor], alvo) {
			i.Seguindo[autor] = append(listaAutor, alvo)
		}
	} else {
		i.Seguindo[autor] = []string{alvo}
	}
	if listaAlvo, fazparteAlvo := i.Seguidores[alvo]; fazparteAlvo {
		if !slices.Contains(i.Seguidores[alvo], autor) {
			i.Seguidores[alvo] = append(listaAlvo, autor)
		}
	} else {
		i.Seguidores[alvo] = []string{autor}
	}
	return nil
}

func (i *Index) InicializaEstado(e *estado.Estado) {
	i.estado = e
}
