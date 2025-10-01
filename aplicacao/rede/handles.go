package rede

import (
	"github.com/freehandle/breeze/crypto"
	"github.com/freehandle/handles/attorney"

	"mega/protocolo/estado"
)

type UserInfo struct {
	Handle  string
	Details string
}

type HandlesDB struct {
	TokenToHandle map[crypto.Token]UserInfo
	HandleToToken map[string]crypto.Token
	Attorneys     map[crypto.Token]struct{}
	AppMEGA       crypto.Token
}

func (a *HandlesDB) Arroba(token crypto.Token) *estado.InfoUsuario {
	if user, ok := a.TokenToHandle[token]; ok {
		return &estado.InfoUsuario{
			Arroba: user.Handle,
		}
	}
	return nil
}

func (a *HandlesDB) Token(handle string) *crypto.Token {
	if token, ok := a.HandleToToken[handle]; ok {
		return &token
	}
	return nil
}

func (a *HandlesDB) IncorporateJoin(action []byte) {
	join := attorney.ParseJoinNetwork(action)
	if join == nil {
		return
	}
	a.TokenToHandle[join.Author] = UserInfo{
		Handle:  join.Handle,
		Details: join.Details,
	}
	a.HandleToToken[join.Handle] = join.Author
}

func (a *HandlesDB) IncorporateUpdate(action []byte) {
	update := attorney.ParseUpdateInfo(action)
	if update == nil {
		return
	}
	handle, ok := a.TokenToHandle[update.Author]
	if !ok {
		return
	}
	a.TokenToHandle[update.Author] = UserInfo{
		Handle:  handle.Handle,
		Details: update.Details,
	}
}

func (a *HandlesDB) IncorporateGrant(action []byte) {
	grant := attorney.ParseGrantPowerOfAttorney(action)
	if grant == nil {
		return
	}
	if grant.Attorney.Equal(a.AppMEGA) {
		a.Attorneys[grant.Author] = struct{}{}
	}
}

func (a *HandlesDB) IncorporateRevoke(action []byte) {
	revoke := attorney.ParseRevokePowerOfAttorney(action)
	if revoke == nil {
		return
	}
	if revoke.Attorney.Equal(a.AppMEGA) {
		delete(a.Attorneys, revoke.Author)
	}
}

func (a *HandlesDB) Incorporate(action []byte) []byte {
	switch attorney.Kind(action) {
	case attorney.VoidType:
		//fmt.Println("MEGA protocol code")
		return action
	case attorney.JoinNetworkType:
		a.IncorporateJoin(action)
		return nil
	case attorney.UpdateInfoType:
		a.IncorporateUpdate(action)
		return nil
	case attorney.GrantPowerOfAttorneyType:
		a.IncorporateGrant(action)
		return nil
	case attorney.RevokePowerOfAttorneyType:
		a.IncorporateRevoke(action)
		return nil
	}
	return action
}
