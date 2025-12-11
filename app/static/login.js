	  // ===============================
	  // 🗂️ "Banco de dados" simulado
	  // ===============================
	  // Aqui criamos um array (lista) de objetos, onde cada objeto representa um usuário
	  // com suas informações de login e dados pessoais.
	  const usuarios = [
		{ nome: "admin", senha: "123456"},
		{ nome: "ana", senha: "ribeiro"},
		{ nome: "tiago", senha: "mittermayer"},
        { nome: "luci", senha: "tabuti"},
        { nome: "gilles", senha: "leite"},
		{ nome: "fabio", senha: "musarra"},
		{ nome: "lari", senha: "lienko"},
	  ];


	  // 🚪 Função para sair do sistema
	  function sair() {
		localStorage.removeItem("usuarioLogado");
		window.location.href = "login.html";
	  }
     
	  // 🔐 Função principal de login
	  function validarLogin() {
		const usuario = document.getElementById("usuario").value.trim();
		const senha = document.getElementById("senha").value.trim();

		// Seleciona o elemento HTML onde será exibida a mensagem de erro ou sucesso
		const msg = document.getElementById("mensagem");

		// Limpa mensagens anteriores antes de validar novamente
		msg.innerText = "";

		
		// 🧠 Validação inicial dos campos
		// Se o usuário ou senha estiverem vazios, exibe uma mensagem de alerta
		if (usuario === "" || senha === "") 
		{
		  msg.innerText = "⚠️ Preencha todos os campos!";
          msg.className = "msgAlerta textoCard corTerciaria"; // aplica cor de alerta à mensagem
		  return; // interrompe a execução da função
		}

		// 🔍 Procura usuário no "banco"
		// Usa o método find() para procurar no array "usuarios" um objeto
		
        // que tenha o mesmo nome e senha informados pelo usuário.
		const usuarioEncontrado = usuarios.find(
		  u => u.nome === usuario && u.senha === senha
		);

		// ===============================
		// ✅ Caso o login seja bem-sucedido
		// ===============================
		if (usuarioEncontrado) 
		{
		  msg.innerText = "✅ Login realizado com sucesso!";
		  msg.className = "msgAlerta textoCard corSucesso"; // aplica cor verde à mensagem

		  // Armazena o usuário logado no navegador (como se fosse uma sessão)
		  // O JSON.stringify transforma o objeto em texto para poder salvar.
		  localStorage.setItem("usuarioLogado", JSON.stringify(usuarioEncontrado));
            
          console.log("Salvando no localStorage:", JSON.stringify(usuarioEncontrado));
          localStorage.setItem("usuarioLogado", JSON.stringify(usuarioEncontrado));

		  // Redireciona o usuário para a página "informacoes.html" após 1,5 segundos
		  setTimeout(() => {
			window.location.href = "meu_jornal.html";
		  }, 1500);
		} 

		// ===============================
		// ❌ Caso o login seja inválido
		// ===============================
		else 
		{
		  msg.innerText = "❌ Usuário ou senha incorretos!";
		  msg.className = "msgAlerta textoCard corErro"; // aplica cor vermelha à mensagem
		}
	  }

      // ===============================
      // 🚪 Função para sair (logout)
      // ===============================
      function logout(){
        // Remove o usuário logado do armazenamento local
        localStorage.removeItem("usuarioLogado");
        
        // Redireciona para a página de login
        window.location.href = "login.html";
      }


	  //verifica usuário para colocar título do site conforme o usuário logado
	  function verificarUsuario() {
		const usuarioLogado = JSON.parse(localStorage.getItem("usuarioLogado"));
		if (usuarioLogado) {
			document.title = `Jornal de ${usuarioLogado.nome}`;
		}
	  }

	  verificarUsuario();

	
	  function nomeiaJornal() {
		const tituloJornal = document.getElementById('tituloJornal');
		const usuarioLogado = JSON.parse(localStorage.getItem("usuarioLogado"));

		if (usuarioLogado) {
			tituloJornal.textContent = `Jornal de ${usuarioLogado.nome}`;
			return;
		}
	}

		document.addEventListener('DOMContentLoaded', () => {
			const tituloJornal = document.getElementById('tituloJornal');
			if (tituloJornal) {
				nomeiaJornal();
			}
		});

		function descubraSenha() {
			alert("Use seu NOME para usuário e seu SOBRENOME para a senha. Utilize apenas letras minúsculas");
		}
