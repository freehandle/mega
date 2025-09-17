package aplicacao

import "time"

type PostarCauso struct {
	Acao       string    `json:"acao"`
	ID         int       `json:"id"`
	CampoCauso string    `json:"campocauso"`
	DataHora   time.Time `json:"dataHora"`
}

type PostarFofoca struct {
	Acao        string    `json:"acao"`
	ID          int       `json:"id"`
	CampoFofoca string    `json:"campofofoca"`
	DataHora    time.Time `json:"dataHora"`
}
