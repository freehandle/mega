package app

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/freehandle/breeze/crypto"
	"github.com/freehandle/mega/indice"
	"github.com/freehandle/mega/protocolo/acoes"
)

var NomesCategorias = [6]string{"meme", "fofoca", "causo", "musica", "ideia", "livro"}

type InformacaoCabecalho struct {
	Arroba string
	// Ativo           string
	Erro      string
	NomeMucua string
	// LinkSelecionada string
}

type ViewConvite struct {
	Cabecalho InformacaoCabecalho
	Seed      string
	// pra testar
	Nome  string
	Nome2 string
}

type ViewPublicar struct {
	Cabecalho InformacaoCabecalho
	Tipo      string
}

// Para puxar do endereco o tipo de pagina a ser construida
type VerPagina struct {
	Usuario   string
	Categoria string
	Data      string
	Tipo      string
}

// Para construir a pagina de um post especifico aberto para leitura post_aberto.html
type PaginaPostAberto struct {
	Categoria    string
	CategoriaMin string
	Arroba       string
	DataPostagem string
	Conteudo     string
	TipoTexto    bool
}

// Para construir o objeto card que vai na pagina jorna.html ou pagina meu_jornal.html
type ConteudoCard struct {
	Categoria       string //categoria com primeira em maiuscula
	CategoriaMin    string //categoria tudo minuscula
	Vazio           bool   //se nao ha postagem da categoria -> true
	Data            string //data da postagem
	ConteudoParcial string // conteudo parcial da postagem (deve caber no card)
}

type Calendario struct {
	DiaAtual       string
	MesAtual       string
	AnoAtual       string
	VermelhosDias  []string
	VermelhosMeses []string
	VermelhosAnos  []string
}

// Para construir a pagina de um jornal sem login jornal.html
type PaginaJornal struct {
	Arroba     string         //dono do jornal
	Cards      []ConteudoCard //vetor com conteudo dos cards
	Calendario Calendario
}

// Para enviar para o montador de card
type ParaMontarCards struct {
	Jornal    *indice.Jornal
	Tipo      string
	Categoria string
	Data      uint64
	Ultimo    int //index do ultimo elemento trazido no caso de posts selecionados via categoria
}

// Encontrando o tipo de pagina a ser construida a partir do endereco acessado
func (v *VerPagina) PegarInfoURL(r *http.Request, mucua string) {
	categorias_possiveis := []string{"causo", "fofoca", "ideia", "livro", "meme", "musica"}
	endereco := r.URL.RequestURI()
	novo := strings.Replace(endereco, mucua, "", 1) // remove servidor

	re := regexp.MustCompile(`\/`) // regex para encontrar o separdor
	partes := re.Split(novo, -1)   // partes separadas

	if len(partes) == 0 || len(partes) > 3 {
		v.Tipo = "erro"
		fmt.Printf("Endereco URL fora de formato")
		return
	}
	if len(partes) == 1 {
		v.Usuario = partes[0]
		v.Tipo = "atual"
		return
	}
	if len(partes) == 2 {
		v.Usuario = partes[0]
		if res, err := regexp.MatchString("^[0-9]{6,6}$", partes[1]); res && err != nil {
			v.Data = partes[1]
			v.Tipo = "data"
			return
		}
		for _, c := range categorias_possiveis {
			if partes[1] == c {
				v.Categoria = partes[1]
				v.Tipo = "categoria"
				return
			}
		}
		return
	}
	if len(partes) == 3 {
		v.Usuario = partes[0]
		v.Categoria = partes[1]
		v.Data = partes[2]
		v.Tipo = "postagem_aberta"
		return
	}
	fmt.Printf("Endereco URL fora de formato")
	return
}

// Cria cards para mostrar no jornal a parti da data e categoria dadas
func (c *ConteudoCard) CriaCard(paraMontar ParaMontarCards) {

	var conteudoTexto *indice.ConteudoData
	var conteudoHash *indice.HashData

	switch paraMontar.Categoria {

	case "ideia":
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
				for _, post := range paraMontar.Jornal.Ideias {
					if post.Data == paraMontar.Data {
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
			c.Data = strconv.FormatUint(conteudoTexto.Data, 10)
			c.ConteudoParcial = conteudoTexto.Conteudo[:200]
			return
		} else {
			c.Vazio = true
			return
		}
	case "causo":
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
				for _, post := range paraMontar.Jornal.Causos {
					if post.Data == paraMontar.Data {
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
			c.Data = strconv.FormatUint(conteudoTexto.Data, 10)
			c.ConteudoParcial = conteudoTexto.Conteudo[:200]
			return
		} else {
			c.Vazio = true
			return
		}
	case "musica":
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
				for _, post := range paraMontar.Jornal.Musicas {
					if post.Data == paraMontar.Data {
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
			c.Data = strconv.FormatUint(conteudoTexto.Data, 10)
			c.ConteudoParcial = conteudoTexto.Conteudo[:200]
			return
		} else {
			c.Vazio = true
			return
		}
	case "fofoca":
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
				for _, post := range paraMontar.Jornal.Fofocas {
					if post.Data == paraMontar.Data {
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
			c.Data = strconv.FormatUint(conteudoTexto.Data, 10)
			c.ConteudoParcial = conteudoTexto.Conteudo[:200]
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
				for _, post := range paraMontar.Jornal.Memes {
					if post.Data == paraMontar.Data {
						conteudoHash = post
						break
					}
				}
				if conteudoTexto == nil {
					c.Vazio = true
					return
				}
			}
			c.Vazio = false
			c.Data = strconv.FormatUint(conteudoHash.Data, 10)
			c.ConteudoParcial = conteudoHash.Hash.String()[0:200]
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
				for _, post := range paraMontar.Jornal.Livros {
					if post.Data == paraMontar.Data {
						conteudoHash = post
						break
					}
				}
				if conteudoTexto == nil {
					c.Vazio = true
					return
				}
			}
			c.Vazio = false
			c.Data = strconv.FormatUint(conteudoHash.Data, 10)
			c.ConteudoParcial = conteudoHash.Hash.String()[0:200]
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

// Gerenciador
func (a *Aplicacao) ManejoInterfacePublicar(w http.ResponseWriter, r *http.Request) {
	arroba := a.Autor(r)
	fmt.Println("ARROBA:", arroba)
	if arroba == "" {
		http.Redirect(w, r, "/credenciais", http.StatusSeeOther)
		return
	}
	view := ViewPublicar{
		Cabecalho: InformacaoCabecalho{
			Arroba:    arroba,
			NomeMucua: a.NomeMucua,
			// Ativo:           "",
			// LinkSelecionada: "",
		},
		Tipo: "causo",
	}
	if err := a.templates.ExecuteTemplate(w, "novotexto.html", view); err != nil {
		log.Println(err)
	}
}

// Gerenciador do template principal da aplicacao
func (a *Aplicacao) ManejoPrincipal(w http.ResponseWriter, r *http.Request) {
	view := InformacaoCabecalho{
		Arroba:    a.Autor(r),
		NomeMucua: a.NomeMucua,
		// Ativo:           "",
		// LinkSelecionada: "",
	}
	if err := a.templates.ExecuteTemplate(w, "main.html", view); err != nil {
		log.Println(err)
	}
}

// Manejo da pagina de jornal sem login
func (a *Aplicacao) ManejoJornal(w http.ResponseWriter, r *http.Request) {

	ver := VerPagina{}
	ver.PegarInfoURL(r, a.NomeMucua)
	if ver.Tipo == "postagem_aberta" || ver.Tipo == "erro" {
		log.Println("endereco nao esta em formato jornal ou contem erro")
		return
	}
	pagina := PaginaJornal{
		Arroba: ver.Usuario,
		Cards:  []ConteudoCard{},
	}

	// tentando pegar o jornal da arroba indicada pelo endereco URL
	jornal, ok := a.Indice.ArrobaParaJornal[ver.Usuario]
	if !ok {
		log.Println("Nome de usuario nao tem jornal associado")
		return
	}

	// marcando a data atual no calendario
	agora := time.Now()
	pagina.Calendario.DiaAtual = agora.Format("02")
	pagina.Calendario.MesAtual = agora.Format("01")
	pagina.Calendario.AnoAtual = agora.Format("2006")

	paraMontar := ParaMontarCards{
		Jornal: jornal,
		Tipo:   ver.Tipo,
	}
	if ver.Categoria != "" {
		paraMontar.Categoria = ver.Categoria
	}
	if ver.Data != "" {
		if dataConvertida, err := strconv.ParseUint(ver.Data, 10, 64); err != nil {
			paraMontar.Data = dataConvertida
		} else {
			log.Println("Erro ao converter data")
		}
	}
	pagina.Cards = CriaCards(paraMontar)
	if err := a.templates.ExecuteTemplate(w, "jornal.html", pagina); err != nil {
		log.Println(err)
	}
}

// Manejo do proprio jornal de usuario (precisa estar logado), tem link para postagem
func (a *Aplicacao) ManejoMeuJornal(w http.ResponseWriter, r *http.Request) {
	pagina := InformacaoCabecalho{
		Arroba:    a.Autor(r),
		NomeMucua: a.NomeMucua,
		// Ativo:           "",
		// LinkSelecionada: "",
	}

	if err := a.templates.ExecuteTemplate(w, "meu_jornal.html", pagina); err != nil {
		log.Println(err)
	}
}

// Manejo de um post do jornal aberto no detalhe
func (a *Aplicacao) ManejoPostAberto(w http.ResponseWriter, r *http.Request) {

	ver := VerPagina{}
	ver.PegarInfoURL(r, a.NomeMucua)

	// checa se o endereco tem a forma de um post especifico
	if ver.Tipo != "postagem_aberta" || ver.Tipo == "erro" {
		log.Println("endereco nao e de uma postagem aberta ou contem erro")
		return
	}
	pagina := PaginaPostAberto{
		Categoria:    ver.Categoria,
		Arroba:       ver.Usuario,
		CategoriaMin: strings.ToLower(ver.Categoria),
	}
	var post_texto *indice.ConteudoData
	var post_hash *indice.HashData

	switch ver.Categoria {
	case "causo":
		posts := a.Indice.ArrobaParaJornal[ver.Usuario].Causos
		for i := 0; i < len(posts); i++ {
			datastr := strconv.FormatUint(posts[i].Data, 10)
			if datastr == ver.Data {
				post_texto = posts[i]
				break
			}
		}
	case "fofoca":
		posts := a.Indice.ArrobaParaJornal[ver.Usuario].Fofocas
		for i := 0; i < len(posts); i++ {
			datastr := strconv.FormatUint(posts[i].Data, 10)
			if datastr == ver.Data {
				post_texto = posts[i]
				break
			}
		}
	case "ideia":
		posts := a.Indice.ArrobaParaJornal[ver.Usuario].Ideias
		for i := 0; i < len(posts); i++ {
			datastr := strconv.FormatUint(posts[i].Data, 10)
			if datastr == ver.Data {
				post_texto = posts[i]
				break
			}
		}
	case "livro":
		posts := a.Indice.ArrobaParaJornal[ver.Usuario].Livros
		for i := 0; i < len(posts); i++ {
			datastr := strconv.FormatUint(posts[i].Data, 10)
			if datastr == ver.Data {
				post_hash = posts[i]
				break
			}
		}
	case "meme":
		posts := a.Indice.ArrobaParaJornal[ver.Usuario].Memes
		for i := 0; i < len(posts); i++ {
			datastr := strconv.FormatUint(posts[i].Data, 10)
			if datastr == ver.Data {
				post_hash = posts[i]
				break
			}
		}
	case "musica":
		posts := a.Indice.ArrobaParaJornal[ver.Usuario].Musicas
		for i := 0; i < len(posts); i++ {
			datastr := strconv.FormatUint(posts[i].Data, 10)
			if datastr == ver.Data {
				post_texto = posts[i]
				break
			}
		}
	default:
		log.Println("categoria nao encontrada")
	}
	if post_texto != nil {
		pagina.DataPostagem = ver.Data
		pagina.Conteudo = post_texto.Conteudo
		pagina.TipoTexto = true
	}
	if post_hash != nil {
		pagina.DataPostagem = ver.Data
		pagina.Conteudo = post_hash.Hash.String()
		pagina.TipoTexto = false
	}
	if err := a.templates.ExecuteTemplate(w, "post_aberto.html", pagina); err != nil {
		log.Println(err)
	}
}

func (a *Aplicacao) ManejoSignin(w http.ResponseWriter, r *http.Request) {
	hashEncoded := r.URL.Path
	hashEncoded = strings.Replace(hashEncoded, "/signin/", "", 1)
	hash := crypto.DecodeHash(hashEncoded)
	fmt.Println("oia", len(a.Convidar))
	if _, ok := a.Convidar[hash]; ok || len(a.Convidar) == 0 {
		view := ViewConvite{
			Cabecalho: InformacaoCabecalho{
				NomeMucua: "",
				// Ativo:           "",
				// LinkSelecionada: "",
				Arroba: "",
			},
			Seed:  hashEncoded,
			Nome:  "teste2",
			Nome2: "teste3",
		}
		fmt.Println("Seed:", hashEncoded)
		if err := a.templates.ExecuteTemplate(w, "signin.html", view); err != nil {
			log.Println(err)
		}
	} else {
		view := InformacaoCabecalho{
			Erro:      "convite inválido",
			NomeMucua: a.NomeMucua,
		}
		if err := a.templates.ExecuteTemplate(w, "login.html", view); err != nil {
			log.Println(err)
		}
	}
}

func (a *Aplicacao) ManejoNovoUsuario(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}
	arroba := r.FormValue("handle")
	email := r.FormValue("email")
	senha := r.FormValue("password")
	ok := a.Gerente.OnboardSigner(arroba, email, senha)
	aviso := InformacaoCabecalho{
		NomeMucua: a.NomeMucua,
	}
	if !ok {
		aviso.Erro = "Confira seu email para ativar sua conta ou tente outro arroba."
	}
	if err := a.templates.ExecuteTemplate(w, "login.html", aviso); err != nil {
		log.Println(err)
	}
	return
}

func (a *Aplicacao) ManejoCredenciais(w http.ResponseWriter, r *http.Request) {
	cookie, arroba, err := a.Gerente.CredentialsHandler(r)
	fmt.Println(cookie)
	fmt.Println(arroba)
	if err != nil {
		aviso := InformacaoCabecalho{
			NomeMucua: a.NomeMucua,
			Erro:      err.Error(),
		}
		if err := a.templates.ExecuteTemplate(w, "login.html", aviso); err != nil {
			log.Println(err)
		}
		return
	}
	// fmt.Println("DEU CERTO AQUI SEU MOCO")
	http.SetCookie(w, cookie)
	/*aviso := InformacaoCabecalho{
		NomeMucua: a.NomeMucua,
		Arroba:    arroba,
	}
	if err := a.templates.ExecuteTemplate(w, "main.html", aviso); err != nil {
		log.Println(err)
	}*/
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return
}

func (a *Aplicacao) ManejoPublica(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}
	handle := a.Autor(r)
	token, ok := a.Indice.ArrobaParaToken[handle]
	if !ok {
		http.Error(w, "usuario desconhecido", http.StatusMethodNotAllowed)
		return
	}
	conteudo := r.FormValue("conteudo")
	tipo := "causo" //r.FormValue("Tipo")
	fmt.Println("TIPO:", tipo)
	if tipo == "causo" {
		causo := &acoes.PostarCauso{
			Epoca:    a.Epoca,
			Autor:    token,
			Conteudo: conteudo,
		}
		if !causo.ValidarFormato() {
			http.Error(w, "formato errado", http.StatusMethodNotAllowed)
			return
		}
		fmt.Println("CAUSO VÁLIDO, ENVIANDO PARA A REDE")
		a.Gateway.Encaminha([]acoes.Acao{causo}, token, a.Epoca)
	}
	aviso := InformacaoCabecalho{
		NomeMucua: a.NomeMucua,
		Arroba:    handle,
	}
	if err := a.templates.ExecuteTemplate(w, "main.html", aviso); err != nil {
		log.Println(err)
	}
}
