package app

import (
	"strconv"
	"time"
)

type DiaCor struct {
	Numero string
	Cor    string // padrao, selecionada, post ou atual -> cinza, azul, cinza escuro ou vermelho
}

type MesAno struct {
	Dias []DiaCor
	Mes  string
	Ano  string
}

type Calendario struct {
	MesAno []MesAno
}

var mapaMeses = map[int]string{1: "Janeiro", 2: "Fevereiro", 3: "Março", 4: "Abril", 5: "Maio", 6: "Junho",
	7: "Julho", 8: "Agosto", 9: "Setembro", 10: "Outubro", 11: "Novembro", 12: "Dezembro"}

// Retorna quantidade de dias no mes
func VetorDiasMes(inicioMes time.Time, datasPostagens []time.Time) []DiaCor {

	// identifica inicio do mes seguinte
	mesSeguinte := inicioMes.AddDate(0, 1, 0)
	// pega o ultimo dia do mes
	ultimoDia := mesSeguinte.Add(-time.Second)

	// pega os dias de postagens daquele mes+ano
	postagens := []int{}
	for _, dia := range datasPostagens {
		if dia.Year() == inicioMes.Year() && dia.Month() == inicioMes.Month() {
			postagens = append(postagens, dia.Day())
		}
	}

	// pega o vetor de dias do mes e marca a cor de cada
	vetorDias := make([]DiaCor, ultimoDia.Day())
	for i := range vetorDias {
		vetorDias[i] = DiaCor{}
		vetorDias[i].Numero = strconv.Itoa(i + 1)
		vetorDias[i].Cor = "padrao"
		for _, d := range postagens {
			if d == i {
				vetorDias[i].Cor = "post"
				break
			}
		}
	}
	return vetorDias
}

// Cria calendarios usados - data 1 é data atual (se possivel), ou data selecionada. data 2, se diferente
//
//	é data selecionada se data 1 for a data atual.
func (c *Calendario) CriaCalendario(data1 time.Time, data1Atual bool, data2 time.Time, datasPostagens []time.Time) {

	c.MesAno = []MesAno{}

	// pega a primeira data do mes 1
	inicioMes1 := time.Date(data1.Year(), data1.Month(), 1, 0, 0, 0, 0, data1.Location())
	mesAno1 := MesAno{}
	mesAno1.Dias = VetorDiasMes(inicioMes1, datasPostagens)
	mesAno1.Mes = mapaMeses[int(data1.Month())] // aparece no titulo do calendario
	mesAno1.Ano = strconv.Itoa(inicioMes1.Year())

	// inclui tag de cor para dia atual
	if data1Atual {
		mesAno1.Dias[data1.Day()+1].Cor = "atual"
	}

	// pega o mes anterior
	mesAnterior := inicioMes1.Add(-time.Second)
	inicioMes2 := time.Date(mesAnterior.Year(), mesAnterior.Month(), 1, 0, 0, 0, 0, mesAnterior.Location())
	mesAno2 := MesAno{}
	mesAno2.Dias = VetorDiasMes(inicioMes2, datasPostagens)
	mesAno2.Mes = mapaMeses[int(inicioMes2.Month())]
	mesAno2.Ano = strconv.Itoa(inicioMes2.Year())

	// inclui tag de cor para data selecionada se houver
	if data1 != data2 {
		if data1.Month() == data2.Month() && data1.Year() == data2.Year() {
			mesAno1.Dias[data2.Day()+1].Cor = "selecionada"
		} else {
			if inicioMes2.Month() == data2.Month() && inicioMes2.Year() == data2.Year() {
				mesAno2.Dias[data2.Day()+1].Cor = "selecionada"
			}
		}
	}
	c.MesAno = append(c.MesAno, mesAno1, mesAno2)
	return

}
