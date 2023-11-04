package main

import (
	"ProjectX/network"
	"fmt"
	"time"
)

func main() {
	trLocal := network.NewLocalTransport("LOCAL")
	trRemote := network.NewLocalTransport("REMOTE")

	trLocal.Connect(trRemote)
	trRemote.Connect(trLocal)

	go func() {
		for {
			msg := []byte("Hello shhhhhh")
			trRemote.SendMessage(trLocal.Addr(), msg)
			time.Sleep(1 * time.Second)
		}
	}()

	opts := network.ServerOpts{
		Transports: []network.Transport{trLocal},
	}

	s := network.NewServer(opts)

	s.Start()
	fmt.Println("dddd")
}
