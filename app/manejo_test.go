package app

import "testing"

func TestVerPagina(t *testing.T) {
	ver := VerPagina{}
	ver.PegarInfoURL("/miga/jornal/ruben", "miga")
	if ver.Tipo != "atual" {
		t.Errorf("Esperado tipo 'atual', obtido '%s'", ver.Tipo)
	}
	ver.PegarInfoURL("/jornal/ruben", "")
	if ver.Tipo != "atual" {
		t.Errorf("Esperado tipo 'atual', obtido '%s'", ver.Tipo)
	}
}

func TestProcessaURL(t *testing.T) {
	pagina := ProcessaURL("/fulano")
	if pagina.Tipo != "atual" {
		t.Errorf("Esperado tipo 'atual', obtido '%s': '%+v'", pagina.Tipo, pagina)
	}
	if pagina.Usuario != "fulano" {
		t.Errorf("Esperado arroba 'fulano', obtido '%s'", pagina.Usuario)
	}
	pagina = ProcessaURL("/fulano/musica")
	if pagina.Tipo != "categoria" {
		t.Errorf("Esperado tipo 'categoria', obtido '%s': '%+v'", pagina.Tipo, pagina)
	}
	if pagina.Usuario != "fulano" {
		t.Errorf("Esperado arroba 'fulano', obtido '%s'", pagina.Usuario)
	}
	pagina = ProcessaURL("/fulano/20250202")
	if pagina.Tipo != "data" {
		t.Errorf("Esperado tipo 'data', obtido '%s': '%+v'", pagina.Tipo, pagina)
	}
	if pagina.Usuario != "fulano" {
		t.Errorf("Esperado arroba 'fulano', obtido '%s'", pagina.Usuario)
	}
	pagina = ProcessaURL("/fulano/20250202/musica")
	if pagina.Tipo != "postagem_aberta" {
		t.Errorf("Esperado tipo 'postagem_aberta', obtido '%s': '%+v'", pagina.Tipo, pagina)
	}
	if pagina.Usuario != "fulano" {
		t.Errorf("Esperado arroba 'fulano', obtido '%s'", pagina.Usuario)
	}
	pagina = ProcessaURL("www.qualquercoisa.com/miga/fulano/musica/20250202/")
	if pagina.Tipo != "postagem_aberta" {
		t.Errorf("Esperado tipo 'postagem_aberta', obtido '%s': '%+v'", pagina.Tipo, pagina)
	}
	if pagina.Usuario != "fulano" {
		t.Errorf("Esperado arroba 'fulano', obtido '%s'", pagina.Usuario)
	}
	pagina = ProcessaURL("/musica/20250202/")
	if pagina.Tipo != "erro" {
		t.Errorf("Esperado tipo 'erro', obtido '%s': '%+v'", pagina.Tipo, pagina)
	}
}
