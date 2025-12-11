// Verifica mês atual
function verificaMesAtual() {

    // pega o mês atual
    const mesAtual = new Date().getMonth();
    const mesAnterior = mesAtual - 1;
    const meses = ['Janeiro', 'Fevereiro', 'Março', 'Abril', 'Maio', 'Junho', 'Julho', 'Agosto', 'Setembro', 'Outubro', 'Novembro', 'Dezembro'];
    const nomeMesAtual = meses[mesAtual];
    const nomeMesAnterior = meses[mesAnterior];
    console.log(nomeMesAtual);
    console.log(nomeMesAnterior);


    // altera o nome do mês atual
    const tituloMesAtual = document.getElementById('tituloMesAtual');
    tituloMesAtual.textContent = `${nomeMesAtual}`;

    // altera o nome do mês anterior
    const tituloMesAnterior = document.getElementById('tituloMesAnterior');
    tituloMesAnterior.textContent = `${nomeMesAnterior}`;
    
}

// verifica o dia atual
function verificaDiaAtual() {
    const diaAtual = new Date().getDate();
    console.log(diaAtual);

    const tituloMesAtual = document.getElementById('tituloMesAtual');

    if (!tituloMesAtual) {
        console.log('Não foi possível encontrar o elemento tituloMesAtual');
        return;
    } 
    else {
        console.log('Elemento tituloMesAtual encontrado: ' + tituloMesAtual.textContent);
    }
    const mesAtual = tituloMesAtual.closest('.mes');
    if (!mesAtual) {
        return;
    } 
    const dias = mesAtual.getElementsByClassName('dia');

    for (let dia of dias) {
        if (diaAtual == parseInt(dia.textContent)) {
            dia.classList.add('diaAtual');
            dia.classList.remove('dia');
        }
    }
}

document.addEventListener('DOMContentLoaded', verificaMesAtual);
document.addEventListener('DOMContentLoaded', verificaDiaAtual);