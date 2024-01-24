package core

import (
	"ProjectX/crypto"
	"ProjectX/types"
	"fmt"
)

type Transaction struct {
	Data []byte

	From crypto.PublicKey
	Signature *crypto.Signature

	hash types.Hash

	firstSeen int64
}

func (t *Transaction) Hash(hasher Hasher[*Transaction]) types.Hash {
	if t.hash.IsZero() {
		t.hash = hasher.Hash(t)
	}
	return t.hash
}

func NewTransaction(data []byte) *Transaction {
	return &Transaction{
		Data: data,
	}
}

func (t *Transaction) Sign(privKey crypto.PrivateKey) error {

	sig, err := privKey.Sign(t.Data)
	if err != nil {
		return err
	}

	t.Signature = sig
	t.From = privKey.PublicKey()
	return nil
}


func (t *Transaction) Verify() error {

	if t.Signature == nil {
		return fmt.Errorf("transaction has no signature")
	}

	if !t.Signature.Verify(t.From, t.Data){
		return fmt.Errorf("transaction has Invalid signature")
	}
	
	return nil
}


func (t *Transaction) Encode(enc Encoder[*Transaction]) error{
	return enc.Encode(t)
}
func (t *Transaction) Decode(dec Decoder[*Transaction]) error{
	return dec.Decode(t)
}

func (t *Transaction) SetFirstSeen(val int64) {
	t.firstSeen = val;
}

func (t *Transaction) FirstSeen() int64 {
	return t.firstSeen;
}