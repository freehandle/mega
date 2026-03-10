package app

import (
	"fmt"
	"log"

	"github.com/freehandle/mega/indice"
)

// Para enviar para o montador de card
type ParaMontarCards struct {
	Jornal    *indice.Jornal
	Tipo      string
	Categoria string
	Data      uint64
	Ultimo    int //index do ultimo elemento trazido no caso de posts selecionados via categoria
	Aplicacao *Aplicacao
}

// Para construir o objeto card que vai na pagina jorna.html
type ConteudoCard struct {
	Categoria       string //categoria com primeira em maiuscula
	CategoriaMin    string //categoria tudo minuscula
	Vazio           bool   //se nao ha postagem da categoria -> true
	Data            string //data da postagem
	DataRef         string // data formatada para ir no endereco de referencia caso clicado
	ConteudoParcial string // conteudo parcial da postagem (deve caber no card)
	Texto           bool   // true se for texto, caso contrario é imagem
}

// Cria cards para mostrar no jornal a parti da data e categoria dadas
func (c *ConteudoCard) CriaCard(paraMontar ParaMontarCards) {

	var conteudoTexto *indice.ConteudoData
	var conteudoHash *indice.HashData

	switch paraMontar.Categoria {

	case "ideia":
		c.Texto = true
		c.Categoria = "Ideia"
		c.CategoriaMin = paraMontar.Categoria
		if len(paraMontar.Jornal.Ideias) > 0 {
			// pegando entrada mais recente de ideia do jornal
			if paraMontar.Tipo == "atual" {
				conteudoTexto = paraMontar.Jornal.Ideias[len(paraMontar.Jornal.Ideias)-1]
			}
			// por categoria, vai pegar o indice pedido
			if paraMontar.Tipo == "categoria" {
				conteudoTexto = paraMontar.Jornal.Ideias[len(paraMontar.Jornal.Ideias)-paraMontar.Ultimo]
			}
			// por data, vai procurar a data pedida
			if paraMontar.Tipo == "data" {
				for n := len(paraMontar.Jornal.Ideias) - 1; n >= 0; n-- {
					post := paraMontar.Jornal.Ideias[n]
					if post.Data <= paraMontar.Data+24*3600 {
						conteudoTexto = post
						break
					}
				}
				if conteudoTexto == nil {
					c.Vazio = true
					return
				}
			}
			c.Vazio = false
			c.Data = DataFormatadaParaCard(paraMontar.Aplicacao, conteudoTexto.Data)
			c.DataRef = DataFormatadaParaReferencia(paraMontar.Aplicacao, conteudoTexto.Data)
			if len(conteudoTexto.Conteudo) > maxLetrasCard {
				c.ConteudoParcial = conteudoTexto.Conteudo[:maxLetrasCard] + "..."
			} else {
				c.ConteudoParcial = conteudoTexto.Conteudo
			}
			return
		} else {
			c.Vazio = true
			return
		}
	case "causo":
		c.Texto = true
		c.Categoria = "Causo"
		c.CategoriaMin = paraMontar.Categoria
		if len(paraMontar.Jornal.Causos) > 0 {
			// pegando entrada mais recente de causo do jornal
			if paraMontar.Tipo == "atual" {
				conteudoTexto = paraMontar.Jornal.Causos[len(paraMontar.Jornal.Causos)-1]
			}
			// por categoria, vai pegar o indice pedido
			if paraMontar.Tipo == "categoria" {
				conteudoTexto = paraMontar.Jornal.Causos[len(paraMontar.Jornal.Causos)-paraMontar.Ultimo]
			}
			// por data, vai procurar a data pedida
			if paraMontar.Tipo == "data" {
				for n := len(paraMontar.Jornal.Causos) - 1; n >= 0; n-- {
					post := paraMontar.Jornal.Causos[n]
					if post.Data <= paraMontar.Data {
						conteudoTexto = post
						break
					}
				}
				if conteudoTexto == nil {
					c.Vazio = true
					return
				}
			}
			c.Vazio = false
			c.Data = DataFormatadaParaCard(paraMontar.Aplicacao, conteudoTexto.Data)
			c.DataRef = DataFormatadaParaReferencia(paraMontar.Aplicacao, conteudoTexto.Data)
			if len(conteudoTexto.Conteudo) > maxLetrasCard {
				c.ConteudoParcial = conteudoTexto.Conteudo[:maxLetrasCard] + "..."
			} else {
				c.ConteudoParcial = conteudoTexto.Conteudo
			}
			return
		} else {
			c.Vazio = true
			return
		}
	case "musica":
		c.Texto = true
		c.Categoria = "Música"
		c.CategoriaMin = paraMontar.Categoria
		if len(paraMontar.Jornal.Musicas) > 0 {
			// pegando entrada mais recente de musica do jornal
			if paraMontar.Tipo == "atual" {
				conteudoTexto = paraMontar.Jornal.Musicas[len(paraMontar.Jornal.Musicas)-1]
			}
			// por categoria, vai pegar o indice pedido
			if paraMontar.Tipo == "categoria" {
				conteudoTexto = paraMontar.Jornal.Musicas[len(paraMontar.Jornal.Musicas)-paraMontar.Ultimo]
			}
			// por data, vai procurar a data pedida
			if paraMontar.Tipo == "data" {
				for n := len(paraMontar.Jornal.Musicas) - 1; n >= 0; n++ {
					post := paraMontar.Jornal.Musicas[n]
					if post.Data <= paraMontar.Data {
						conteudoTexto = post
						break
					}
				}
				if conteudoTexto == nil {
					c.Vazio = true
					return
				}
			}
			c.Vazio = false
			c.Data = DataFormatadaParaCard(paraMontar.Aplicacao, conteudoTexto.Data)
			c.DataRef = DataFormatadaParaReferencia(paraMontar.Aplicacao, conteudoTexto.Data)
			if len(conteudoTexto.Conteudo) > maxLetrasCard {
				c.ConteudoParcial = conteudoTexto.Conteudo[:maxLetrasCard] + "..."
			} else {
				c.ConteudoParcial = conteudoTexto.Conteudo
			}
			return
		} else {
			c.Vazio = true
			return
		}
	case "fofoca":
		c.Texto = true
		c.Categoria = "Fofoca"
		c.CategoriaMin = paraMontar.Categoria
		if len(paraMontar.Jornal.Fofocas) > 0 {
			// pegando entrada mais recente de fofoca do jornal
			if paraMontar.Tipo == "atual" {
				conteudoTexto = paraMontar.Jornal.Fofocas[len(paraMontar.Jornal.Fofocas)-1]
			}
			// por categoria, vai pegar o indice pedido
			if paraMontar.Tipo == "categoria" {
				conteudoTexto = paraMontar.Jornal.Fofocas[len(paraMontar.Jornal.Fofocas)-paraMontar.Ultimo]
			}
			// por data, vai procurar a data pedida
			if paraMontar.Tipo == "data" {
				for n := len(paraMontar.Jornal.Fofocas) - 1; n >= 0; n-- {
					post := paraMontar.Jornal.Fofocas[n]
					if post.Data <= paraMontar.Data {
						conteudoTexto = post
						break
					}
				}
				if conteudoTexto == nil {
					c.Vazio = true
					return
				}
			}
			c.Vazio = false
			c.Data = DataFormatadaParaCard(paraMontar.Aplicacao, conteudoTexto.Data)
			c.DataRef = DataFormatadaParaReferencia(paraMontar.Aplicacao, conteudoTexto.Data)
			if len(conteudoTexto.Conteudo) > maxLetrasCard {
				c.ConteudoParcial = conteudoTexto.Conteudo[:maxLetrasCard] + "..."
			} else {
				c.ConteudoParcial = conteudoTexto.Conteudo
			}
			return
		} else {
			c.Vazio = true
			return
		}
	case "meme":
		c.Categoria = "Meme"
		c.CategoriaMin = paraMontar.Categoria
		if len(paraMontar.Jornal.Memes) > 0 {
			// pegando entrada mais recente de meme do jornal
			if paraMontar.Tipo == "atual" {
				conteudoHash = paraMontar.Jornal.Memes[len(paraMontar.Jornal.Memes)-1]
			}
			// por categoria, vai pegar o indice pedido
			if paraMontar.Tipo == "categoria" {
				conteudoHash = paraMontar.Jornal.Memes[len(paraMontar.Jornal.Memes)-paraMontar.Ultimo]
			}
			// por data, vai procurar a data pedida
			if paraMontar.Tipo == "data" {
				for n := len(paraMontar.Jornal.Memes) - 1; n >= 0; n-- {
					post := paraMontar.Jornal.Memes[n]
					if post.Data <= paraMontar.Data {
						conteudoHash = post
						break
					}
				}
				if conteudoHash == nil {
					c.Vazio = true
					return
				}
			}
			c.Vazio = false
			c.Data = DataFormatadaParaCard(paraMontar.Aplicacao, conteudoHash.Data)
			c.DataRef = DataFormatadaParaReferencia(paraMontar.Aplicacao, conteudoHash.Data)

			c.ConteudoParcial = fmt.Sprintf("%s%s", conteudoHash.Hash.String(), conteudoHash.Tipo)

			return
		} else {
			c.Vazio = true
			return
		}
	case "livro":
		c.Categoria = "Livro"
		c.CategoriaMin = paraMontar.Categoria
		if len(paraMontar.Jornal.Livros) > 0 {
			// pegando entrada mais recente de livro do jornal
			if paraMontar.Tipo == "atual" {
				conteudoHash = paraMontar.Jornal.Livros[len(paraMontar.Jornal.Livros)-1]
			}
			// por categoria, vai pegar o indice pedido
			if paraMontar.Tipo == "categoria" {
				conteudoHash = paraMontar.Jornal.Livros[len(paraMontar.Jornal.Livros)-paraMontar.Ultimo]
			}
			// por data, vai procurar a data pedida
			if paraMontar.Tipo == "data" {
				for n := len(paraMontar.Jornal.Livros) - 1; n >= 0; n-- {
					post := paraMontar.Jornal.Livros[n]
					if post.Data <= paraMontar.Data {
						conteudoHash = post
						break
					}
				}
				if conteudoHash == nil {
					c.Vazio = true
					return
				}
			}
			c.Vazio = false
			c.Data = DataFormatadaParaCard(paraMontar.Aplicacao, conteudoHash.Data)
			c.DataRef = DataFormatadaParaReferencia(paraMontar.Aplicacao, conteudoHash.Data)
			c.ConteudoParcial = fmt.Sprintf("%s%s", conteudoHash.Hash.String(), conteudoHash.Tipo)
			return
		} else {
			c.Vazio = true
			return
		}
	default:
		c.Vazio = true
		log.Println("erro ao criar card")
		return
	}
}

// Monta o array de cards para a pagina jornal
func CriaCards(paraMontar ParaMontarCards) []ConteudoCard {
	// categoria precisa estar com a mesma grafia (no plural) usada para a construcao da struct index.Jornal

	var vetorCards []ConteudoCard

	// sem especificacao de categoria ou data -> traz os mais atuais por categoria
	if paraMontar.Tipo == "atual" {
		for _, cat := range NomesCategorias {
			paraMontar.Categoria = cat
			card := ConteudoCard{}
			card.CriaCard(paraMontar)
			vetorCards = append(vetorCards, card)
		}
		return vetorCards
	}
	if paraMontar.Tipo == "categoria" {
		var total int = 0
		switch paraMontar.Categoria {
		case "ideia":
			total = len(paraMontar.Jornal.Ideias)
		case "meme":
			total = len(paraMontar.Jornal.Memes)
		case "fofoca":
			total = len(paraMontar.Jornal.Fofocas)
		case "causo":
			total = len(paraMontar.Jornal.Causos)
		case "musica":
			total = len(paraMontar.Jornal.Musicas)
		case "livro":
			total = len(paraMontar.Jornal.Livros)
		}
		if total > 0 {
			i := 0
			for i < total {
				card := ConteudoCard{}
				paraMontar.Ultimo = i
				card.CriaCard(paraMontar)
				vetorCards = append(vetorCards, card)
			}
		}
		return vetorCards
		// v := reflect.ValueOf(jornal)
		// posts := v.FieldByName(categoria)
	}
	if paraMontar.Tipo == "data" {
		for _, cat := range NomesCategorias {
			paraMontar.Categoria = cat
			card := ConteudoCard{}
			card.CriaCard(paraMontar)
			vetorCards = append(vetorCards, card)
		}
	}
	return vetorCards
}
