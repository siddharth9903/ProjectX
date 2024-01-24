package core

import (
	"ProjectX/crypto"
	"ProjectX/types"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func randomBlock(height uint32, prevBlockHash types.Hash) *Block {
	header := &Header{
		Version: 1,
		PrevBlockHash: prevBlockHash,
		Timestamp: time.Now().UnixNano(),
		Height: height,
	}


	return NewBlock(header, []Transaction{})
}

func randomBlockWithSignature(t *testing.T, height uint32, prevBlockHash types.Hash) *Block {
	privKey := crypto.GeneratePrivateKey()

	b := randomBlock(height, prevBlockHash)
	tx := randomTxWithSignature(t)
	b.AddTransaction(tx)
	assert.Nil(t, b.Sign(privKey))

	return b
}

func TestHashBlock(t *testing.T) {
	b := randomBlock(0,types.Hash{})
	fmt.Println("b.Hash", b.Hash(BlockHasher{}))
}

func TestSignBlock(t *testing.T) {
	privKey := crypto.GeneratePrivateKey()

	b := randomBlock(1, types.Hash{})
	assert.Nil(t, b.Sign(privKey))
	assert.NotNil(t, b.Signature)
	assert.Equal(t, b.Validator, privKey.PublicKey())
}

func TestVerifyBlock(t *testing.T) {

	privKey := crypto.GeneratePrivateKey()
	anotherPrivKey := crypto.GeneratePrivateKey()

	b := randomBlock(1, types.Hash{})

	assert.Nil(t, b.Sign(privKey))
	b.Signature = nil
	assert.NotNil(t, b.Verify())

	assert.Nil(t, b.Sign(privKey))
	b.Validator = anotherPrivKey.PublicKey()
	assert.NotNil(t, b.Verify())
}

// import (
// 	"ProjectX/types"
// 	"bytes"
// 	"fmt"
// 	"testing"
// 	"time"

// 	"github.com/stretchr/testify/assert"
// )

// func TestHeader_Encode_Decode(t *testing.T){
// 	h := &Header{
// 		Version:1,
// 		PrevBlock: types.RandomHash(),
// 		Timestamp: time.Now().UnixNano(),
// 		Height:32,
// 		Nonce:12450,
// 	}

// 	buf := &bytes.Buffer{}
// 	assert.Nil(t, h.EncodeBinary(buf))

// 	hDecode := &Header{}
// 	assert.Nil(t, hDecode.DecodeBinary(buf))
// 	assert.Equal(t, hDecode, h)
// }

// func TestBlock_Encode_Decode(t *testing.T){
// 	b := &Block{
// 		Header: Header{
// 		Version:1,
// 		PrevBlock: types.RandomHash(),
// 		Timestamp: time.Now().UnixNano(),
// 		Height:32,
// 		Nonce:12450,
// 		},
// 		Transactions: nil,
// 	}

// 	buf := &bytes.Buffer{}
// 	assert.Nil(t, b.EncodeBinary(buf))

// 	bDecode := &Block{}
// 	assert.Nil(t, bDecode.DecodeBinary(buf))
// 	assert.Equal(t, bDecode, b)
// }

// func TestBlockHash(t *testing.T){
// 	b := &Block{
// 		Header: Header{
// 		Version:1,
// 		PrevBlock: types.RandomHash(),
// 		Timestamp: time.Now().UnixNano(),
// 		Height:32,
// 		Nonce:12450,
// 		},
// 		Transactions: []Transaction{},
// 	}

// 	h := b.Hash()
// 	assert.False(t, h.IsZero())

// 	fmt.Printf("%+v",h)
// }