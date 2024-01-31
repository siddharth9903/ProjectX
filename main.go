package main

import (
	"ProjectX/core"
	"ProjectX/crypto"
	"ProjectX/network"
	"bytes"
	"math/rand"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

func main() {
	trLocal := network.NewLocalTransport("LOCAL")
	trRemote := network.NewLocalTransport("REMOTE")

	trLocal.Connect(trRemote)
	trRemote.Connect(trLocal)

	go func() {
		for {
			// msg := []byte("Hello Local")
			// trRemote.SendMessage(trLocal.Addr(), msg)
			// sendTransaction(trLocal)
			if err := sendTransaction(trRemote, trLocal.Addr()); err != nil {
				logrus.Error(err)
			}
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

func sendTransaction(tr network.Transport, to network.NetAddr) error {
	Pk := crypto.GeneratePrivateKey()
	tx := core.NewTransaction([]byte(strconv.FormatInt(int64(rand.Intn(100)),10)))
	tx.Sign(Pk)

	txBytes := &bytes.Buffer{}
	err := tx.Encode(core.NewGobTxEncoder(txBytes))
	if err != nil {
		return err
	}

	msg := network.NewMessage(network.MessageTypeTx, txBytes.Bytes())
	
	return tr.SendMessage(to, msg.Bytes())
}
