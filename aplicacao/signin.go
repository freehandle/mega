package aplicacao

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"mega/aplicacao/configuracoes"
	"mega/protocolo/acoes"
	"net/smtp"
	"time"

	"github.com/freehandle/breeze/crypto"
	"github.com/freehandle/safe"
)

const wellcomeBody = `Seu endereço de email foi associado ao handle %s.  @s livres, protocolo mega`

const resetBody = `Para resetar clique em https://%v`

// var emailSigninMessage = "To: %v\r\n" + "Subject: Protocolo MEGA\r\n" + "\r\n" + "%v\r\n"

const signinWithoutHandleBody = `Your email was associated to a MEGA account for the handle %v. Acressar %v using the fingerprint %v `

const signinBody = `handle %v. power of attorney to %v using the fingerprint %v`

// var emailPasswordMessage = "To: %v\r\n" + "Subject: Senha MEGA \r\n" + "\r\n" + "%v\r\n"

const passwordMessage = `Your new password for app mega account for the handle %v is %v`

type GerenciadorSignin struct {
	cofre             *safe.Safe // for optional direct onboarding
	pendente          []*Signerin
	senhas            configuracoes.PasswordManager
	TokenDoProcurador crypto.Token
	mail              Mailer
	Procurador        *ProcuradorGeral
	Procuracao        map[string]crypto.Token
}

type SMTPGmail struct {
	Password string
	From     string
}

func (s *SMTPGmail) Send(to, subject, body string) bool {
	auth := smtp.PlainAuth("", s.From, s.Password, "smtp.gmail.com")
	emailMsg := fmt.Sprintf("To: %s\r\n"+"Subject: %s\r\n"+"\r\n"+"%s\r\n", to, subject, body)
	err := smtp.SendMail("smtp.gmail.com:587", auth, s.From, []string{to}, []byte(emailMsg))
	if err != nil {
		log.Printf("email sending error: %v", err)
		return false
	}
	return true
}

type Mailer interface {
	Send(to, subject, body string) bool
}

type Signerin struct {
	Handle      string
	Email       string
	TimeStamp   time.Time
	FingerPrint string
}

func NewSigninManager(passwords configuracoes.PasswordManager, mail Mailer, attorney *ProcuradorGeral) *GerenciadorSignin {
	if attorney == nil {
		log.Print("PANIC BUG: NewSigninManager called with nil attorney ")
		return nil
	}
	return &GerenciadorSignin{
		cofre:             attorney.safe,
		pendente:          make([]*Signerin, 0),
		senhas:            passwords,
		TokenDoProcurador: attorney.Token,
		Procurador:        attorney,
		mail:              mail,
		Procuracao:        make(map[string]crypto.Token),
	}
}

func (s *GerenciadorSignin) RequestReset(user crypto.Token, email, domain string) bool {
	if !s.senhas.HasWithEmail(user, email) {
		return false
	}
	reset := s.senhas.AddReset(user, email)
	url := fmt.Sprintf("%s/r/%s", domain, reset)
	if reset == "" {
		return false
	}
	go s.mail.Send(email, "Reset de senha pra protocolo MEGA", fmt.Sprintf(resetBody, url))
	return true
}

func (s *GerenciadorSignin) Reset(user crypto.Token, url, password string) bool {
	return s.senhas.DropReset(user, url, password)
}

func (s *GerenciadorSignin) Check(user crypto.Token, password string) bool {
	hashed := crypto.Hasher(append(user[:], []byte(password)...))
	return s.senhas.Check(user, hashed)

}

func (s *GerenciadorSignin) Set(user crypto.Token, password string, email string) {
	hashed := crypto.Hasher(append(user[:], []byte(password)...))
	s.senhas.Set(user, hashed, email)
}

func (s *GerenciadorSignin) DirectReset(user crypto.Token, newpassword string) bool {
	newhashed := crypto.Hasher(append(user[:], []byte(newpassword)...))
	return s.senhas.Reset(user, newhashed)
}

func (s *GerenciadorSignin) Has(token crypto.Token) bool {
	return s.senhas.Has(token)
}

func (s *GerenciadorSignin) OnboardSigner(handle, email, passwd string) bool {
	if s.cofre == nil {
		log.Println("PANIC BUG: OnboardSigner called with nil safe")
		return false
	}
	ok, token := s.cofre.SigninWithToken(handle, passwd, email)
	if !ok {
		return false
	}
	if err := s.cofre.GrantPower(handle, s.TokenDoProcurador.Hex(), crypto.EncodeHash(crypto.HashToken(token))); err != nil {
		log.Println("error granting power of attorney", err)
		return false
	}
	s.Set(token, passwd, email)
	signin := acoes.Entrar{
		Epoca:   s.Procurador.estado.Epoca,
		Autor:   token,
		Reasons: "Sign in no protocolo MEGA aprovado com procuração",
	}
	s.Procurador.Send([]acoes.Acao{&signin}, token)
	s.Procuracao[handle] = token
	go s.mail.Send(email, "Boas vindas ao protocolo MEGA", fmt.Sprintf(wellcomeBody, handle, handle, handle, handle, handle))
	return true
}

func (s *GerenciadorSignin) AddSigner(handle, email string, token *crypto.Token) {
	signer := &Signerin{}
	for _, pending := range s.pendente {
		if signer.Handle == handle {
			signer = pending
		}
	}
	signer.Handle = handle
	signer.Email = email
	t, _ := crypto.RandomAsymetricKey()
	signer.FingerPrint = crypto.EncodeHash(crypto.HashToken(t))
	signer.TimeStamp = time.Now()
	if token != nil {
		go s.sendSigninEmail(signinBody, handle, email, signer.FingerPrint)
	} else {
		go s.sendSigninEmail(signinWithoutHandleBody, handle, email, signer.FingerPrint)
	}
	s.pendente = append(s.pendente, signer)
}

func randomPassword() string {
	bytes := make([]byte, 10)
	rand.Read(bytes)
	return base64.StdEncoding.EncodeToString(bytes)
}

func (s *GerenciadorSignin) DarProcuracao(token crypto.Token, arroba, fingerprint string) {
	for n, signer := range s.pendente {
		if signer.Handle == arroba && signer.FingerPrint == fingerprint {
			passwd := randomPassword()
			s.Set(token, passwd, signer.Email)
			signin := acoes.Entrar{
				Epoca:   s.Procurador.estado.Epoca,
				Autor:   token,
				Reasons: "Ingresso na aplicação MEGA com procuração",
			}
			s.Procurador.Send([]acoes.Acao{&signin}, token)
			go s.sendPasswordEmail(signer.Handle, signer.Email, passwd)
			s.pendente = append(s.pendente[:n], s.pendente[n+1:]...)
			s.Procuracao[arroba] = token
		}
	}
}

func (s *GerenciadorSignin) RevogarProcuracao(token crypto.Token, arroba string, assinatura string) {
	// for n, signer := range s.pendente {
	// 	if signer.Handle == arroba &&  == assinatura.String() {
	// 		passwd := randomPassword()
	// 		s.Set(token, passwd, signer.Email)
	// 		signin := acoes.Entrar{
	// 			Epoca:   s.Procurador.estado.Epoca,
	// 			Autor:   token,
	// 			Reasons: "Ingresso na aplicação MEGA com procuração",
	// 		}
	// 		s.Procurador.Send([]acoes.Acao{&signin}, token)
	// 		go s.sendPasswordEmail(signer.Handle, signer.Email, passwd)
	// 		s.pendente = append(s.pendente[:n], s.pendente[n+1:]...)
	// 		s.Procuracao[arroba] = token
	// 	}
	// }
}

func (s *GerenciadorSignin) sendSigninEmail(msg, handle, email, fingerprint string) {
	body := fmt.Sprintf(msg, handle, s.TokenDoProcurador, fingerprint)
	s.mail.Send(email, "Entrando no protocolo MEGA", body)
}

func (s *GerenciadorSignin) sendPasswordEmail(handle, email, password string) {
	body := fmt.Sprintf(passwordMessage, handle, password)
	s.mail.Send(email, "Senha do MEGA", body)
}
