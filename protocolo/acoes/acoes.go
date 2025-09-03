package acoes

import "github.com/freehandle/breeze/crypto"

// associando cada acao prevista pelo protocolo a um byte identificador
const (
	APostarCauso byte = iota
	APostarFofoca
	APostarIdeia
	APostarLivro
	APostarMeme
	APostarMusica
	APostarErro // pro caso de dar um erro no byte da acao
)

// cada acao declarada apos AIdeia tem como identificador o byte da ação anterior + 1
type Acao interface {
	Serializa() []byte     // serializa a instrucao
	Autoria() crypto.Token // verifica a autoria
	FazHash() crypto.Hash  // faz o hash da instrucao de postagem
	// ChecaTempo() bool      // checa se pode postar, ou se ainda nao pode
	// ChecaFormato() bool    // checa se o formato está adequado para o tipo de acao
}

// retorna o byte do tipo de acao referente aos dados recebidos
func TipoDeAcao(dados []byte) byte {
	if len(dados) < 8+crypto.TokenSize+1 {
		return APostarErro
	}
	byteAcao := dados[8+crypto.TokenSize]
	if byteAcao >= APostarErro {
		return APostarErro
	}
	return byteAcao
}
