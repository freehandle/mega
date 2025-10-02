package rede

import (
	"log"

	"github.com/freehandle/mega/aplicacao"

	"github.com/freehandle/breeze/util"
	"github.com/freehandle/handles/attorney"
)

type Gateway interface {
	SendAction(action []byte) error
}

type MEGANode struct {
	Gateway Gateway
	Axe     *HandlesDB
	General *aplicacao.ProcuradorGeral
}

// 0 = new block
// 1 = action?

type Signal struct {
	Signal byte
	Data   []byte
}

func ByteArrayToSignal(receive chan []byte) chan *Signal {
	signals := make(chan *Signal)
	go func() {
		for {
			action, ok := <-receive
			if !ok {
				return
			}
			if len(action) > 0 {
				signals <- &Signal{Signal: action[0], Data: action[1:]}
			}
		}
	}()
	return signals
}

// canal é um primitivo de sincronia
// canal <- value manda para o canal
// value = <-canal recebe do canal
func NewMEGANode(axe *HandlesDB, attorneyGeneral *aplicacao.ProcuradorGeral, signals chan *Signal) {
	for {
		signal := <-signals
		if signal.Signal == 0 {
			epoch, _ := util.ParseUint64(signal.Data, 0)
			attorneyGeneral.DefineEpoca(epoch)
		} else if signal.Signal == 1 {
			if attorney.IsAxeNonVoid(signal.Data) {
				if attorney.Kind(signal.Data) == attorney.GrantPowerOfAttorneyType {
					grant := attorney.ParseGrantPowerOfAttorney(signal.Data)
					if grant != nil {
						if attorneyGeneral.Token.Equal(grant.Attorney) {
							if user, ok := axe.TokenToHandle[grant.Author]; ok {
								attorneyGeneral.IncorporarProcuracao(user.Handle, grant)
							}
						}
					}
				} else if attorney.Kind(signal.Data) == attorney.RevokePowerOfAttorneyType {
					revoke := attorney.ParseRevokePowerOfAttorney(signal.Data)
					if revoke != nil {
						if attorneyGeneral.Token.Equal(revoke.Attorney) {
							if user, ok := axe.TokenToHandle[revoke.Author]; ok {
								attorneyGeneral.IncorporarRevogacao(user.Handle, revoke)
							}
						}
					}
				} else if attorney.Kind(signal.Data) == attorney.JoinNetworkType {
					join := attorney.ParseJoinNetwork(signal.Data)
					if join != nil {
						axe.IncorporateJoin(signal.Data)
					}
				} else if attorney.Kind(signal.Data) == attorney.UpdateInfoType {
					axe.IncorporateUpdate(signal.Data)
				}
			}
			acaoMEGA := axe.Incorporate(signal.Data)
			if acaoMEGA != nil {
				action := BreezeToMEGA(signal.Data)
				attorneyGeneral.Incorporar(action)
			}
		} else {
			log.Printf("invalid signal: %v", signal.Signal)
		}
	}
}
