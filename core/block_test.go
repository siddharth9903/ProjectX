package core

import (
	"ProjectX/crypto"
	"ProjectX/types"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func randomBlock(t *testing.T, height uint32, prevBlockHash types.Hash) *Block {
	privKey := crypto.GeneratePrivateKey()
	tx := randomTxWithSignature(t)

	header := &Header{
		Version: 1,
		PrevBlockHash: prevBlockHash,
		Timestamp: time.Now().UnixNano(),
		Height: height,
	}

	b, err := NewBlock(header, []*Transaction{tx})
	assert.Nil(t, err)

	dataHash,err := CalculateDataHash([]*Transaction{tx})
	assert.Nil(t, err)
	header.DataHash = dataHash

	assert.Nil(t, b.Sign(privKey))
	return b
}

func TestHashBlock(t *testing.T) {
	b := randomBlock(t,0,types.Hash{})
	fmt.Println("b.Hash", b.Hash(BlockHasher{}))
}

func TestSignBlock(t *testing.T) {
	privKey := crypto.GeneratePrivateKey()

	b := randomBlock(t, 1, types.Hash{})
	assert.Nil(t, b.Sign(privKey))
	assert.NotNil(t, b.Signature)
	assert.Equal(t, b.Validator, privKey.PublicKey())
}

func TestVerifyBlock(t *testing.T) {

	privKey := crypto.GeneratePrivateKey()
	anotherPrivKey := crypto.GeneratePrivateKey()

	b := randomBlock(t, 1, types.Hash{})

	assert.Nil(t, b.Sign(privKey))
	b.Signature = nil
	assert.NotNil(t, b.Verify())

	assert.Nil(t, b.Sign(privKey))
	b.Validator = anotherPrivKey.PublicKey()
	assert.NotNil(t, b.Verify())
}