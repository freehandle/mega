

function compartilharJornal () {
    const url = window.location.href;

    navigator.clipboard.writeText(url);
    alert('URL copiada para a área de transferência');
}