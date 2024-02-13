package main

import (
	"ProjectX/core"
	"ProjectX/crypto"
	"ProjectX/network"
	"bytes"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

func main() {
	trLocal := network.NewLocalTransport("LOCAL")
	trRemoteA := network.NewLocalTransport("REMOTE_A")
	trRemoteB := network.NewLocalTransport("REMOTE_B")
	trRemoteC := network.NewLocalTransport("REMOTE_C")

	trLocal.Connect(trRemoteA)
	trRemoteA.Connect(trRemoteB)
	trRemoteB.Connect(trRemoteC)

	trRemoteA.Connect(trLocal)


	initRemoteServers([]network.Transport{trRemoteA, trRemoteB, trRemoteC})

	go func() {
		for {
			if err := sendTransaction(trRemoteA, trLocal.Addr()); err != nil {
				logrus.Error(err)
			}
			time.Sleep(1 * time.Second)
		}
	}()


	// go func() {
	// 	time.Sleep(7 * time.Second)
		
	// 	trRemoteLate := network.NewLocalTransport("LATE_REMOTE")
	// 	trRemoteC.Connect(trRemoteLate)
		
	// 	lateServer := makeServer("LATE_REMOTE",trRemoteLate, nil)
	// 	go lateServer.Start()
	// }()

	privKey := crypto.GeneratePrivateKey()
	localServer := makeServer("LOCAL", trLocal, &privKey)

	localServer.Start()
}

func initRemoteServers(transports []network.Transport){
	for i, tr := range transports {
		id := fmt.Sprintf("REMOTE_%d",i)
		s := makeServer(id,tr, nil)
		go s.Start()
	}
}

func makeServer(id string, tr network.Transport, privKey *crypto.PrivateKey) *network.Server {
	opts := network.ServerOpts{
		Transports: []network.Transport{tr},
		PrivateKey : privKey,
		ID : id,
	}

	s, err := network.NewServer(opts)
	if err != nil {
		logrus.Error(err)
	}

	return s
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
