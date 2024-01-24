package main

import (
	"ProjectX/crypto"
	"ProjectX/network"
	"time"
)

func main() {
	trLocal := network.NewLocalTransport("LOCAL")
	trRemote := network.NewLocalTransport("REMOTE")

	trLocal.Connect(trRemote)
	trRemote.Connect(trLocal)

	go func() {
		for {
			msg := []byte("Hello Local")
			trRemote.SendMessage(trLocal.Addr(), msg)
			time.Sleep(1 * time.Second)
		}
	}()

	Pk := crypto.GeneratePrivateKey()
	opts := network.ServerOpts{
		Transports: []network.Transport{trLocal},
		PrivateKey : &Pk,
	}

	s := network.NewServer(opts)

	s.Start()
}
