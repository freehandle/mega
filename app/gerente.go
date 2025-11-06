package app

import (
	"fmt"
	"iu/auth"

	"github.com/freehandle/breeze/crypto"
)

func ContrataGerente(app *Aplicacao, path, senhaGmail, usuarioGmail string, credenciais crypto.PrivateKey) (*auth.SigninManager, error) {

	arqCofrinho := fmt.Sprintf("%s/%s", path, "senhas.dat")
	cofrinho := auth.NewFilePasswordManager(arqCofrinho)

	mensagens := auth.MessagesTemplates{
		Reset:                    "To reset your password, click the following link: %s",
		ResetHeader:              "Password Reset",
		Signin:                   "To sign in, click the following link: %s",
		SigninHeader:             "Sign In",
		Wellcome:                 "Welcome! Your account has been created.",
		WellcomeHeader:           "Welcome to our App",
		EmailSigninMessage:       "To sign in without your handle, click the following link: %s",
		EmailSigninMessageHeader: "Sign In Without Handle",
		PasswordMessage:          "Your new password is: %s",
		PasswordMessageHeader:    "New Password",
	}

	var gmail auth.Mailer
	if usuarioGmail == "" {
		gmail = auth.TesteGmail{}
	} else {
		gmail = &auth.SMTPGmail{
			Password: senhaGmail,
			From:     usuarioGmail,
		}
	}

	carteiro := &auth.SMTPManager{
		Mail:      gmail,
		Token:     credenciais.PublicKey(),
		Templates: mensagens,
	}

	arqDoceria := fmt.Sprintf("%s/%s", path, "cookies.dat")
	doceria, err := auth.OpenCokieStore(arqDoceria)
	if err != nil {
		return nil, err
	}

	gerente := &auth.SigninManager{
		Passwords:      cofrinho,
		Cookies:        doceria,
		Mail:           carteiro,
		Granted:        make(map[string]crypto.Token),
		Credentials:    credenciais,
		Members:        app,
		SafeAddress:    "http://localhost:8089",
		SafeAPIAddress: "http://localhost:8090",
	}
	return gerente, nil
}
