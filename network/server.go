package network

import (
	"ProjectX/core"
	"ProjectX/crypto"
	"bytes"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

type ServerOpts struct {
	Transports []Transport
	PrivateKey *crypto.PrivateKey
	BlockTime time.Duration
	RPCDecodeFunc RPCDecodeFunc
	RPCProcessor RPCProcessor
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

	if opts.RPCDecodeFunc == nil {
		opts.RPCDecodeFunc = DefaultRPCDecodeFunc
	}

	s := &Server{
		ServerOpts: opts,
		blockTime: opts.BlockTime,
		isValidator: opts.PrivateKey != nil,
		memPool : NewTxPool(),
		rpcChan:    make(chan RPC),
		quitChan:   make(chan struct{}, 1),
	}

	if s.RPCProcessor == nil {
		s.RPCProcessor = s
	}

	return s
}

func (s *Server) Start() {
	s.initTransports()
	ticker := time.NewTicker(s.BlockTime)

free:
	for {
		select {
		case rpc := <-s.rpcChan:
			decodedMessage, err := s.RPCDecodeFunc(rpc)
			if err != nil {
				logrus.Error(err)
			}
			if err := s.RPCProcessor.ProcessMessage(decodedMessage); err!=nil {
				logrus.Error(err)
			}

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

func (s *Server) ProcessMessage (decodedMsg *DecodedMessage) error{
	
	switch t := decodedMsg.Data.(type) {
		case *core.Transaction:
			return s.processTransaction(t)
	}

	return nil
}

func (s *Server) processTransaction(tx *core.Transaction) error {

	hash := tx.Hash(core.TxHasher{})

	if s.memPool.Has(hash) {
		logrus.WithFields(
			logrus.Fields{
				"hash": hash,
			}).Info("transaction already in mempool")
	}

	if err := tx.Verify(); err != nil {
		return err
	}

	tx.SetFirstSeen(time.Now().UnixNano())

	logrus.WithFields(logrus.Fields{
		"hash": hash,
		"mempool length": s.memPool.Len(),
	}).Info("adding new tx to the mempool")

	go s.broadcastTx(tx)
	//add tx to other peers

	return s.memPool.Add(tx)
}

func (s *Server)broadcast(b []byte) error{
	for _, t := range s.Transports{
		t.Broadcast(b)
	}
	return nil
}

func (s *Server) broadcastTx(tx *core.Transaction) error{
	buf := &bytes.Buffer{}
	err := tx.Encode(core.NewGobTxEncoder(buf))

	if err != nil{
		return err
	}

	msg := NewMessage(MessageTypeTx, buf.Bytes())
	return s.broadcast(msg.Bytes())
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
