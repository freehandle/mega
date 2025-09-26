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

const signinWithoutHandleBody = `Your email was associated to a Synergy account for the handle %v. Acressar %v using the fingerprint %v `

const signinBody = `handle %v. power of attorney to %v using the fingerprint %v`

// var emailPasswordMessage = "To: %v\r\n" + "Subject: Senha MEGA \r\n" + "\r\n" + "%v\r\n"

const passwordMessage = `Your new password for app mega account for the handle %v is %v`

type GerenciadorSignin struct {
	safe          *safe.Safe // for optional direct onboarding
	pending       []*Signerin
	passwords     configuracoes.PasswordManager
	AttorneyToken crypto.Token
	mail          Mailer
	Attorney      *ProcuradorGeral
	Granted       map[string]crypto.Token
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
		safe:          attorney.safe,
		pending:       make([]*Signerin, 0),
		passwords:     passwords,
		AttorneyToken: attorney.Token,
		Attorney:      attorney,
		mail:          mail,
		Granted:       make(map[string]crypto.Token),
	}
}

type SigninManager struct {
	safe          *safe.Safe // for optional direct onboarding
	pending       []*Signerin
	passwords     configuracoes.PasswordManager
	AttorneyToken crypto.Token
	mail          Mailer
	Attorney      *ProcuradorGeral
	Granted       map[string]crypto.Token
}

func (s *SigninManager) RequestReset(user crypto.Token, email, domain string) bool {
	if !s.passwords.HasWithEmail(user, email) {
		return false
	}
	reset := s.passwords.AddReset(user, email)
	url := fmt.Sprintf("%s/r/%s", domain, reset)
	if reset == "" {
		return false
	}
	go s.mail.Send(email, "Synergy password reset", fmt.Sprintf(resetBody, url))
	return true
}

func (s *SigninManager) Reset(user crypto.Token, url, password string) bool {
	return s.passwords.DropReset(user, url, password)
}

func (s *SigninManager) Check(user crypto.Token, password string) bool {
	hashed := crypto.Hasher(append(user[:], []byte(password)...))
	return s.passwords.Check(user, hashed)

}

func (s *SigninManager) Set(user crypto.Token, password string, email string) {
	hashed := crypto.Hasher(append(user[:], []byte(password)...))
	s.passwords.Set(user, hashed, email)
}

func (s *SigninManager) DirectReset(user crypto.Token, newpassword string) bool {
	newhashed := crypto.Hasher(append(user[:], []byte(newpassword)...))
	return s.passwords.Reset(user, newhashed)
}

func (s *SigninManager) Has(token crypto.Token) bool {
	return s.passwords.Has(token)
}

func (s *SigninManager) OnboardSigner(handle, email, passwd string) bool {
	if s.safe == nil {
		log.Println("PANIC BUG: OnboardSigner called with nil safe")
		return false
	}
	ok, token := s.safe.SigninWithToken(handle, passwd, email)
	if !ok {
		return false
	}
	if err := s.safe.GrantPower(handle, s.AttorneyToken.Hex(), crypto.EncodeHash(crypto.HashToken(token))); err != nil {
		log.Println("error granting power of attorney", err)
		return false
	}
	s.Set(token, passwd, email)
	signin := acoes.Entrar{
		Epoca:   s.Attorney.estado.Epoca,
		Autor:   token,
		Reasons: "Synergy app sign in with approved power of attorney",
	}
	s.Attorney.Send([]acoes.Acao{&signin}, token)
	s.Granted[handle] = token
	go s.mail.Send(email, "Synergy Protocol Welcome", fmt.Sprintf(wellcomeBody, handle, handle, handle, handle, handle))
	return true
}

func (s *SigninManager) AddSigner(handle, email string, token *crypto.Token) {
	signer := &Signerin{}
	for _, pending := range s.pending {
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
	s.pending = append(s.pending, signer)
}

func randomPassword() string {
	bytes := make([]byte, 10)
	rand.Read(bytes)
	return base64.StdEncoding.EncodeToString(bytes)
}

func (s *SigninManager) RevokeAttorney(handle string) {
	delete(s.Granted, handle)
}

func (s *SigninManager) GrantAttorney(token crypto.Token, handle, fingerprint string) {
	log.Println("to aqui", handle, fingerprint)
	for n, signer := range s.pending {
		if signer.Handle == handle && signer.FingerPrint == fingerprint {
			passwd := randomPassword()
			s.Set(token, passwd, signer.Email)
			signin := acoes.Entrar{
				Epoca:   s.Attorney.estado.Epoca,
				Autor:   token,
				Reasons: "Synergy app sign in with approved power of attorney",
			}
			s.Attorney.Send([]acoes.Acao{&signin}, token)
			go s.sendPasswordEmail(signer.Handle, signer.Email, passwd)
			s.pending = append(s.pending[:n], s.pending[n+1:]...)
			s.Granted[handle] = token
		}
	}
}

func (s *SigninManager) sendSigninEmail(msg, handle, email, fingerprint string) {
	body := fmt.Sprintf(msg, handle, s.AttorneyToken, fingerprint)
	s.mail.Send(email, "Entrando no protocolo MEGA", body)
}

func (s *SigninManager) sendPasswordEmail(handle, email, password string) {
	body := fmt.Sprintf(passwordMessage, handle, password)
	s.mail.Send(email, "Senha do MEGA", body)
}
