package network

import (
	"ProjectX/core"
	"ProjectX/crypto"
	"ProjectX/types"
	"bytes"
	"os"
	"time"

	"github.com/go-kit/log"
)

type ServerOpts struct {
	Transports []Transport
	PrivateKey *crypto.PrivateKey
	BlockTime time.Duration
	RPCDecodeFunc RPCDecodeFunc
	RPCProcessor RPCProcessor
	Logger log.Logger
	ID string
}

const DefaultBlockTime time.Duration= 5 * time.Second;

type Server struct {
	ServerOpts
	blockTime time.Duration
	isValidator bool
	memPool *TxPool
	chain *core.Blockchain
	rpcChan  chan RPC
	quitChan chan struct{}
}

func NewServer(opts ServerOpts) (*Server, error) {
	if opts.BlockTime == time.Duration(0) {
		opts.BlockTime = DefaultBlockTime
	}

	if opts.RPCDecodeFunc == nil {
		opts.RPCDecodeFunc = DefaultRPCDecodeFunc
	}

	if opts.Logger == nil {
		opts.Logger = log.NewLogfmtLogger(os.Stderr)
		opts.Logger = log.With(opts.Logger, "ID", opts.ID)
	}

	chain, err := core.NewBlockchain( opts.Logger,genesisBlock())
	if err != nil {
		return nil, err
	}

	s := &Server{
		ServerOpts: opts,
		chain : chain,
		blockTime: opts.BlockTime,
		isValidator: opts.PrivateKey != nil,
		memPool : NewTxPool(),
		rpcChan:    make(chan RPC),
		quitChan:   make(chan struct{}, 1),
	}

	if s.RPCProcessor == nil {
		s.RPCProcessor = s
	}

	if s.isValidator {
		go s.validatorLoop()
	}

	return s,nil
}

func (s *Server) Start() {
	s.initTransports()

free:
	for {
		select {
		case rpc := <-s.rpcChan:
			decodedMessage, err := s.RPCDecodeFunc(rpc)
			if err != nil {
				s.Logger.Log("Error",err)
			}
			if err := s.RPCProcessor.ProcessMessage(decodedMessage); err!=nil {
				s.Logger.Log("Error",err)
			}
		case <-s.quitChan:
			break free
	
		}
	}

	s.Logger.Log("msg","Server shut down")
}

func (s *Server) validatorLoop() {
	ticker := time.NewTicker(s.blockTime)

	s.Logger.Log("msg", "Starting validator loop", "blocktime", s.blockTime)

	for {
		<- ticker.C
		s.createNewBlock()	
	}
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
		return nil
	}

	if err := tx.Verify(); err != nil {
		return err
	}

	tx.SetFirstSeen(time.Now().UnixNano())

	s.Logger.Log("msg","adding new tx to the mempool","hash",hash,"mempool length", s.memPool.Len())

	go s.broadcastTx(tx)

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

	header,err := s.chain.GetHeader(s.chain.Height())
	if err != nil{
		return err
	}

	txx := s.memPool.Transactions()

	block, err := core.NewBlockFromPrevHeader(header,txx)
	if err != nil{
		return err
	}

	if err := block.Sign(*s.PrivateKey); err != nil{
		return err
	}

	if err := s.chain.AddBlock(block); err != nil{
		return err
	}

	s.memPool.Flush()

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

func genesisBlock() *core.Block {
	h := &core.Header{
			Version:1,
			DataHash: types.Hash{},
			Timestamp: 00000,
			PrevBlockHash:types.Hash{},
			Height:0,
	}

	b,_ := core.NewBlock(h, nil)
	return b
}