    //define categorias 
    const categorias = ['fofoca', 'causo', 'ideia', 'livro', 'meme', 'musica'];

    // Função para carregar as postagens do localStorage
function carregarPostagens() {
    //recupera conteúdos do localStorage
        return JSON.parse(localStorage.getItem('conteudos')) || [];
      }

// Função auxiliar para formatar a data no formato "dia de mês de ano"
function formatarData() {
  const diaAtual = new Date().getDate();
  const anoAtual = new Date().getFullYear();
  const mesAtual = new Date().getMonth();
  // const meses = ['Janeiro', 'Fevereiro', 'Março', 'Abril', 'Maio', 'Junho', 'Julho', 'Agosto', 'Setembro', 'Outubro', 'Novembro', 'Dezembro'];
  // const nomeMes = meses[mesAtual];
  return `${diaAtual}/${mesAtual}/${anoAtual}`;
}


// Função para criar o card de cada postagem segundo a categoria
function criarCard(conteudo) {
    const card = document.createElement('article');
    card.className = 'card sombraG';
    card.setAttribute('data-categoria', conteudo.tipo);
    card.innerHTML = `
      <div class="cabecalhoCard corPrimaria">
        <img class="imgP" src="img/${conteudo.tipo}.png" alt="icone ${conteudo.tipo}">
        <p class="tituloCabecalho">${conteudo.tipo.charAt(0).toUpperCase() + conteudo.tipo.slice(1)}</p>
      </div>
      <p class="textoCard dataCard" id="dataCard"> ${conteudo.data}</p>
      <br>
      <p class="textoCard textoLimitado">${conteudo.texto}</p>
      <br>
      <div class="rodapeCard">
        <img class="imgG" src="img/zoom.png" alt="zoom">
      </div>
    `;
    return card;
  }

  //montar cards vazios
  function criarCardVazio(tipo) {
    const card = document.createElement('article');
    card.className = 'cardVazio sombraG';
    card.setAttribute('data-categoria', tipo);
    card.innerHTML = `
      <div class="cabecalhoCard corPrimaria">
        <img class="imgP" src="img/${tipo}.png" alt="icone ${tipo}">
        <p class="tituloCabecalho">${tipo.charAt(0).toUpperCase() + tipo.slice(1)}</p>
      </div>
      <br>
      <p class="textoCard textoLimitado">Este Jornal não tem ${tipo === 'ideia' ? 'Ideias' : `${tipo.charAt(0).toUpperCase() + tipo.slice(1)}s`} no dia selecionado.</p>
       <div class="rodapeCard">
        <img class="imgG" src="img/zoom.png" alt="zoom">
      </div>
    `;
    return card;
  }

//exibe cards na página
function exibirPostagens() {
    const conteudos = carregarPostagens();
  
    // Monta um objeto com arrays vazios para cada categoria
    const agrupados = categorias.reduce((acc, categoria) => {
      acc[categoria] = [];
      return acc;
    }, {});
  
    // Distribui os conteúdos lidos em seus grupos
    conteudos.forEach(conteudo => {
      if (agrupados[conteudo.tipo]) {
        agrupados[conteudo.tipo].push(conteudo);
      }
    });
  
    // Para cada categoria, busca a seção da página e popula
    categorias.forEach(categoria => {
      const container = document.querySelector(`.listaCards[data-categoria="${categoria}"]`);
      if (!container) return;
  
      container.innerHTML = '';
  
      const posts = agrupados[categoria];
      if (posts.length === 0) {
        // Sem conteúdo salvo: exibe card vazio
        container.appendChild(criarCardVazio(categoria));
      } else {
        // Há conteúdo: cria um card para cada registro
        posts.forEach(conteudo => container.appendChild(criarCard(conteudo)));
      }
    });
  }

  //carrega postagens ao carregar a página
  document.addEventListener('DOMContentLoaded', exibirPostagens);