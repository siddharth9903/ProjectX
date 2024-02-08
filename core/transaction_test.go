package core

import (
	"ProjectX/crypto"
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignTransaction(t *testing.T) {

	privKey := crypto.GeneratePrivateKey()
	tx := &Transaction{
		Data: []byte("tx data 1"),
	}

	assert.Nil(t, tx.Sign(privKey))

	assert.NotNil(t, tx.Signature)
	assert.Equal(t, tx.From, privKey.PublicKey())
}

func TestVerifyTransaction(t *testing.T) {

	privKey := crypto.GeneratePrivateKey()
	anotherPrivKey := crypto.GeneratePrivateKey()

	tx := &Transaction{
		Data: []byte("tx data 1"),
	}

	assert.Nil(t, tx.Sign(privKey))
	tx.Signature = nil
	assert.NotNil(t, tx.Verify())

	assert.Nil(t, tx.Sign(privKey))
	tx.From = anotherPrivKey.PublicKey()
	assert.NotNil(t, tx.Verify())
}

func randomTxWithSignature(t *testing.T) *Transaction { 
	privKey := crypto.GeneratePrivateKey()

	tx := Transaction{
		Data: []byte("tx data"),
	}
	assert.Nil(t, tx.Sign(privKey));
	return &tx;
}

func TestEncodeDecodeTransaction(t *testing.T){
	tx := randomTxWithSignature(t)
	buf := &bytes.Buffer{}
	assert.Nil(t, tx.Encode(NewGobTxEncoder(buf)))

	txDecoded := new(Transaction)
	assert.Nil(t, txDecoded.Decode(NewGobTxDecoder(buf)))
	assert.Equal(t, tx, txDecoded)
}