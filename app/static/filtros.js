//seleciona categorias e filtra cards correspondentes
function selecionarAba(elemento){
      const jaSelecionada = elemento.classList.contains("categoriaSelecionada");

  // Se já estiver selecionada, desmarca e mostra todos os cards
  if (jaSelecionada) {
    elemento.classList.remove("categoriaSelecionada");
    elemento.classList.add("categoria");

    const cards = document.querySelectorAll("[data-categoria]");
    for (let card of cards) {
      card.classList.remove("esconder");
    }
    return;
  }

  // Troca a classe da categoria selecionada
  const abas = document.getElementsByClassName("categoriaSelecionada");
  for (let aba of abas) {
    aba.classList.remove("categoriaSelecionada");
    aba.classList.add("categoria");
  }

  elemento.classList.remove("categoria");
  elemento.classList.add("categoriaSelecionada");

  // Filtra os cards pelo atributo data-categoria
  const categoriaSelecionada = elemento.id;
  const cards = document.querySelectorAll("[data-categoria]");

  for (let card of cards) {
    const categoriaCard = card.getAttribute("data-categoria");
    if (categoriaSelecionada === "todas" || categoriaCard === categoriaSelecionada) {
      card.classList.remove("esconder");
    } else {
      card.classList.add("esconder");
    }
  }
}


//Filtro de postagens segundo data de postagem
function carregarPostagens() {
  //recupera conteúdos do localStorage
      return JSON.parse(localStorage.getItem('conteudos')) || [];
    }

//função para selecionar posts de um dia específico e destacar o dia no calendário
function selecionarDia(elemento){
  const diaSelecionado = elemento.classList.contains("dia");
  if (diaSelecionado) {
    // Remove seleção de outros dias
    const todosDias = document.querySelectorAll(".diaCheio");
    todosDias.forEach(dia => {
      dia.classList.remove("diaCheio");
      dia.classList.add("dia");
    });
    
    elemento.classList.remove("dia");
    elemento.classList.add("diaCheio");
    const diaNumero = parseInt(elemento.textContent);
    localStorage.setItem('diaSelecionado', diaNumero);
    
    filtrarPostagensPorDia();
    return;
  }
  else {
    elemento.classList.remove("diaCheio");
    elemento.classList.add("dia");
    localStorage.removeItem('diaSelecionado'); // Remove do localStorage
    
    // Mostra todas as postagens novamente
    exibirPostagens();
  }
}

//função para filtrar postagens por dia
function filtrarPostagensPorDia() {
  const conteudos = carregarPostagens();
  const diaSelecionado = parseInt(localStorage.getItem('diaSelecionado')); // Converte para número
  
  // Obtém a data atual para comparar mês e ano também
  const hoje = new Date();
  const mesAtual = hoje.getMonth();
  const anoAtual = hoje.getFullYear();
  
  // Filtra postagens do dia selecionado
  const postagensFiltradas = conteudos.filter(conteudo => {
    // Converte a string de data (ex: "12/1/2024") para objeto Date
    // toLocaleDateString() pode retornar formatos diferentes, então precisamos parsear
    const partesData = conteudo.data.split('/');
    const dataPost = new Date(partesData[2], partesData[1] - 1, partesData[0]); // Ano, Mês (0-indexed), Dia
    
    // Compara dia, mês e ano
    return dataPost.getDate() === diaSelecionado && 
           dataPost.getMonth() === mesAtual && 
           dataPost.getFullYear() === anoAtual;
  });
  
  // Agrupa por categoria (mesma lógica de exibirPostagens)
  const categorias = ['fofoca', 'causo', 'ideia', 'livro', 'meme', 'musica'];
  const agrupados = categorias.reduce((acc, categoria) => {
    acc[categoria] = [];
    return acc;
  }, {});
  
  postagensFiltradas.forEach(conteudo => {
    if (agrupados[conteudo.tipo]) {
      agrupados[conteudo.tipo].push(conteudo);
    }
  });
  
  // Exibe apenas as postagens filtradas
  categorias.forEach(categoria => {
    const container = document.querySelector(`.listaCards[data-categoria="${categoria}"]`);
    if (!container) return;
    
    container.innerHTML = '';
    
    const posts = agrupados[categoria];
    if (posts.length > 0) {
      // Cria cards apenas para postagens do dia selecionado
      posts.forEach(conteudo => container.appendChild(criarCard(conteudo)));
    }
    // Se não houver postagens do dia, criar card vazio
    else {
      container.appendChild(criarCardVazio(categoria));
    }
  });
}