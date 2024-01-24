package core

import (
	"fmt"
	"sync"
)

type Blockchain struct {
	store Storage
	lock sync.RWMutex
	headers []*Header
	validator Validator
}

func NewBlockchain(genesis *Block) (*Blockchain, error) {
	bc := &Blockchain{
		store : NewMemoryStore(),
		headers : []*Header{},
	}
	bc.validator = NewBlockValidator(bc)
	bc.addBlockWithoutValidation(genesis)

	return bc, nil
}

func (bc *Blockchain) AddBlock(b *Block) error {
	
	if err:= bc.validator.ValidateBlock(b); err != nil {
		return err
	}

	return bc.addBlockWithoutValidation(b)
}

func (bc *Blockchain) GetHeader(height uint32) (*Header, error){

	if height > bc.Height(){
		return nil, fmt.Errorf("given height %d is too high", height)
	}
	bc.lock.Lock()
	defer bc.lock.Unlock()

	return bc.headers[height], nil
}

func (bc *Blockchain) addBlockWithoutValidation(b *Block) error {
	bc.lock.Lock()
	defer bc.lock.Unlock()
	bc.headers = append(bc.headers,b.Header)

	// logrus.WithFields(logrus.Fields{
	// 	"height": b.Height, 
	// 	"hash": b.Hash(BlockHasher{}),
	// }).Info("adding new block")

	return bc.store.Put(b)
}

func (bc *Blockchain) SetValidator(v Validator) {
	bc.validator = v
}

func (bc *Blockchain) HasBlock(height uint32) bool {
	return bc.Height() >= height
}

func (bc *Blockchain) Height() uint32 {
	bc.lock.RLock()
	defer bc.lock.RUnlock()
	return uint32(len(bc.headers) - 1)
}