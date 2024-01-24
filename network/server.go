package network

import (
	"ProjectX/core"
	"ProjectX/crypto"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

type ServerOpts struct {
	Transports []Transport
	PrivateKey *crypto.PrivateKey
	BlockTime time.Duration
}

const DefaultBlockTime time.Duration= 5 * time.Second;

type Server struct {
	ServerOpts
	blockTime time.Duration
	isValidator bool
	memPool *TxPool
	rpcChan  chan RPC
	quitChan chan struct{}
}

func NewServer(opts ServerOpts) *Server {
	if opts.BlockTime == time.Duration(0) {
		opts.BlockTime = DefaultBlockTime
	}

	return &Server{
		ServerOpts: opts,
		blockTime: opts.BlockTime,
		isValidator: opts.PrivateKey != nil,
		memPool : NewTxPool(),
		rpcChan:    make(chan RPC),
		quitChan:   make(chan struct{}, 1),
	}
}

func (s *Server) Start() {
	s.initTransports()
	ticker := time.NewTicker(s.BlockTime)

free:
	for {
		select {
		case rpc := <-s.rpcChan:
			fmt.Println(rpc)
		case <-s.quitChan:
			break free
		case <-ticker.C:
			if s.isValidator {
				s.createNewBlock()
			}
		}
	}

	fmt.Println("Server shut down")
}

func (s *Server) handleTransaction(tx *core.Transaction) error {
	if err := tx.Verify(); err != nil {
		return err
	}

	hash := tx.Hash(core.TxHasher{})
	if s.memPool.Has(hash) {
		logrus.WithFields(
			logrus.Fields{
				"hash": hash,
			}).Info("transaction already in mempool")
	}

	logrus.WithFields(logrus.Fields{
		"hash": hash,
	}).Info("adding new tx to the mempool")

	return s.memPool.Add(tx)
}


func (s *Server) createNewBlock() error {
	fmt.Println("creating a new block")
	return nil
}


func (s *Server) initTransports() {
	for _, tr := range s.Transports {
		go func(tr Transport) {
			for rpc := range tr.Consume() {
				s.rpcChan <- rpc
			}
		}(tr)
	}
}
