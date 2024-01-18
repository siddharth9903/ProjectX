package core

import (
	"ProjectX/crypto"
	"fmt"
)

type Transaction struct {
	Data []byte

	From crypto.PublicKey
	Signature *crypto.Signature
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

