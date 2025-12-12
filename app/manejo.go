package app

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/freehandle/breeze/crypto"
	"github.com/freehandle/mega/indice"
	"github.com/freehandle/mega/protocolo/acoes"
)

type VerPost struct {
	Usuario      string
	Categoria    string
	DataPostagem string
}

type InformacaoCabecalho struct {
	Arroba string
	// Ativo           string
	Erro      string
	NomeMucua string
	// LinkSelecionada string
}

type ViewPostAberto struct {
	Categoria    string
	CategoriaMin string
	Arroba       string
	DataPostagem string
	Conteudo     string
	TipoTexto    bool
	TipoImagem   bool
}

type ConteudoCard struct {
	Categoria       string //categoria com primeira em maiuscula
	CategoriaMin    string //categoria tudo minuscula
	Vazio           bool   //se nao ha postagem da categoria -> true
	DataPostagem    string //data da postagem
	ConteudoParcial string // conteudo parcial da postagem (deve caber no card)
}

type ViewJornal struct {
	Arroba string         //dono do jornal
	Cards  []ConteudoCard //vetor com conteudo dos cards
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

// Gerenciador do template principal da aplicacao
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

// Manejo de um jornal de usuario
func (a *Aplicacao) ManejoJornal(w http.ResponseWriter, r *http.Request) {
	view := ViewJornal{
		// Arroba: ver,
	}
	if err := a.templates.ExecuteTemplate(w, "jornal.html", view); err != nil {
		log.Println(err)
	}
}

// Manejo do proprio jornal de usuario (precisa estar logado), tem link para postagem
func (a *Aplicacao) ManejoMeuJornal(w http.ResponseWriter, r *http.Request) {
	view := InformacaoCabecalho{
		Arroba:    a.Autor(r),
		NomeMucua: a.NomeMucua,
		// Ativo:           "",
		// LinkSelecionada: "",
	}

	if err := a.templates.ExecuteTemplate(w, "meu_jornal.html", view); err != nil {
		log.Println(err)
	}
}

func (v *VerPost) PegarInfoURL(r *http.Request, mucua string) {
	endereco := r.URL.RequestURI()
	novo := strings.Replace(endereco, mucua, "", 1) // remove servidor

	re := regexp.MustCompile(`\/`) // regex para encontrar o separdor
	partes := re.Split(novo, -1)   // partes separadas

	if len(partes) == 0 {
		return
	}
	if len(partes) == 1 {
		v.Usuario = partes[0]
		return
	}
	if len(partes) == 2 {
		v.Usuario = partes[0]
		v.Categoria = partes[1]
		return
	}
	if len(partes) == 3 {
		v.Usuario = partes[0]
		v.Categoria = partes[1]
		v.DataPostagem = string(partes[2])
		return
	}
	fmt.Printf("Endereco URL fora de formato")
	return
}

// Manejo de um post do jornal aberto no detalhe
func (a *Aplicacao) ManejoPostAberto(w http.ResponseWriter, r *http.Request) {

	ver := VerPost{}
	ver.PegarInfoURL(r, a.NomeMucua)
	view := ViewPostAberto{
		Categoria:    ver.Categoria,
		Arroba:       ver.Usuario,
		CategoriaMin: strings.ToLower(ver.Categoria),
	}
	var post_texto *indice.ConteudoData
	var post_hash *indice.HashData
	switch ver.Categoria {
	case "Causo":
		posts := a.Indice.ArrobaParaJornal[ver.Usuario].Causos
		for i := 0; i < len(posts); i++ {
			datastr := strconv.FormatUint(posts[i].Data, 10)
			if datastr == ver.DataPostagem {
				post_texto = posts[i]
				break
			}
		}
	case "Fofoca":
		posts := a.Indice.ArrobaParaJornal[ver.Usuario].Fofocas
		for i := 0; i < len(posts); i++ {
			datastr := strconv.FormatUint(posts[i].Data, 10)
			if datastr == ver.DataPostagem {
				post_texto = posts[i]
				break
			}
		}
	case "Ideia":
		posts := a.Indice.ArrobaParaJornal[ver.Usuario].Ideias
		for i := 0; i < len(posts); i++ {
			datastr := strconv.FormatUint(posts[i].Data, 10)
			if datastr == ver.DataPostagem {
				post_texto = posts[i]
				break
			}
		}
	case "Livro":
		posts := a.Indice.ArrobaParaJornal[ver.Usuario].Livros
		for i := 0; i < len(posts); i++ {
			datastr := strconv.FormatUint(posts[i].Data, 10)
			if datastr == ver.DataPostagem {
				post_hash = posts[i]
				break
			}
		}
	case "Meme":
		posts := a.Indice.ArrobaParaJornal[ver.Usuario].Memes
		for i := 0; i < len(posts); i++ {
			datastr := strconv.FormatUint(posts[i].Data, 10)
			if datastr == ver.DataPostagem {
				post_hash = posts[i]
				break
			}
		}
	case "Musica":
		posts := a.Indice.ArrobaParaJornal[ver.Usuario].Musicas
		for i := 0; i < len(posts); i++ {
			datastr := strconv.FormatUint(posts[i].Data, 10)
			if datastr == ver.DataPostagem {
				post_texto = posts[i]
				break
			}
		}
	default:
		log.Println("categoria nao encontrada")
	}
	if post_texto != nil {
		view.DataPostagem = ver.DataPostagem
		view.Conteudo = post_texto.Conteudo
		view.TipoTexto = true
	}
	if post_hash != nil {
		view.DataPostagem = ver.DataPostagem
		view.Conteudo = post_hash.Hash.String()
		view.TipoImagem = true
	}
	if err := a.templates.ExecuteTemplate(w, "post_aberto.html", view); err != nil {
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
