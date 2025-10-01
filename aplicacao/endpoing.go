package aplicacao

type InformacaoCabecalho struct {
	ArrobaUsuario string
	Ativo         string
	Erro          string
	NomeMucua     string
}

type Mucua struct {
	Cabecalho InformacaoCabecalho
	Nome      string
}
