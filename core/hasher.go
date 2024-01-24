package core

import (
	"ProjectX/types"
	"crypto/sha256"
)

type Hasher[T any] interface{
	Hash(T) types.Hash
}

type BlockHasher struct{}


func (BlockHasher) Hash(h *Header) types.Hash { 
	b:= sha256.Sum256(h.Bytes())
	return types.Hash(b)
}

type TxHasher struct{}

func (TxHasher) Hash(tx *Transaction) types.Hash {
	return types.Hash(sha256.Sum256(tx.Data))
}
