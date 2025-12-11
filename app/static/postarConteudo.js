
// Função para voltar à página anterior
function voltarPagina() {
  window.history.back();
}

// Função para postar o conteúdo
function postarConteudo() {
  // Obter o valor do textarea
  const conteudo = document.querySelector('textarea').value;
  
    // Verificar se o conteúdo não está vazio
    if (conteudo.trim() === '') {
        alert('Por favor, adicione algum conteúdo antes de postar.');
        return;
    } 
    
    //método para salvar o conteúdo no localStorage
    const tipoConteudo = document.getElementById('tipoConteudo').value;
    const novoConteudo = {
        tipo: tipoConteudo,
        texto: conteudo,
        data: new Date().toLocaleDateString()
    };

    // Recuperar conteúdos existentes do localStorage
    const conteudosExistentes = JSON.parse(localStorage.getItem('conteudos')) || [];
    conteudosExistentes.push(novoConteudo);

    // Salvar o novo conteúdo no localStorage
    localStorage.setItem('conteudos', JSON.stringify(conteudosExistentes));

    // Confirmar postagem e redirecionar para a Home
    alert('Conteúdo postado com sucesso!');
    document.getElementById('textoConteudo').value = '';
    window.location.href = 'home.html';
}
    //faz o botão mudar o rótulo conforme a opção selecionada
function atualizarLabelBotao() {
    const selectTipo = document.getElementById('tipoConteudo');
    const botaoPostar = document.getElementById('botaoPostar');

    if (!selectTipo || !botaoPostar) {
        return;
    }

    const opcaoSelecionada = selectTipo.selectedOptions[0] || selectTipo.options[selectTipo.selectedIndex];
    if (!opcaoSelecionada) {
        return;
    }

    const rotulo = (opcaoSelecionada.dataset.label || opcaoSelecionada.textContent || '').trim();

    if (rotulo) {
        botaoPostar.textContent = ` adicionar ${rotulo} `;
    }
}

document.addEventListener('DOMContentLoaded', () => {
    const selectTipo = document.getElementById('tipoConteudo');

    if (!selectTipo) {
        return;
    }

    atualizarLabelBotao();
    selectTipo.addEventListener('change', atualizarLabelBotao);
});