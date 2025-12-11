// Busca no localStorage o item "usuarioLogado" (salvo como texto JSON)
	  // e converte de volta para objeto JavaScript com JSON.parse()
	  const usuarioLogado = JSON.parse(localStorage.getItem("usuarioLogado"));

	  // 👀 Verifica se há usuário logado
	  if (!usuarioLogado) {
		// Se não houver login registrado, exibe alerta e retorna à página de login
		alert("⚠️ Você precisa fazer login primeiro!");
		window.location.href = "login.html"; // redireciona
	  }