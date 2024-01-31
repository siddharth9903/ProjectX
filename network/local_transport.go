package network

import (
	"bytes"
	"fmt"
	"sync"
)

type LocalTransport struct {
	addr      NetAddr
	consumeCh chan RPC
	lock      sync.RWMutex
	peers     map[NetAddr]*LocalTransport
}

func NewLocalTransport(addr NetAddr) Transport {
	return &LocalTransport{
		addr:      addr,
		consumeCh: make(chan RPC, 1024),
		peers:     make(map[NetAddr]*LocalTransport),
	}
}

func (t *LocalTransport) Consume() <-chan RPC {
	return t.consumeCh
}

func (t *LocalTransport) Connect(tr Transport) error {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.peers[tr.Addr()] = tr.(*LocalTransport)
	return nil
}

func (t *LocalTransport) SendMessage(to NetAddr, message []byte) error {
	t.lock.RLock()
	defer t.lock.RUnlock()

	peer, ok := t.peers[to]

	if !ok {
		return fmt.Errorf("%v failed to send message to %v", t.addr, to)
	}

	peer.consumeCh <- RPC{
		From:    t.addr,
		Payload: bytes.NewReader(message),
	}

	return nil
}

func (t *LocalTransport) Broadcast(payload []byte) error {
	for _, p := range t.peers{
		err := t.SendMessage(p.Addr(), payload)
		if err != nil {
			return err
		}
	}
	return nil
}


func (t *LocalTransport) Addr() NetAddr {
	return t.addr
}
