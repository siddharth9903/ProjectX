package core

import (
	"ProjectX/types"
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHeader_Encode_Decode(t *testing.T){
	h := &Header{
		Version:1,
		PrevBlock: types.RandomHash(),
		Timestamp: time.Now().UnixNano(),
		Height:32,
		Nonce:12450,
	}

	buf := &bytes.Buffer{}
	assert.Nil(t, h.EncodeBinary(buf))

	hDecode := &Header{}
	assert.Nil(t, hDecode.DecodeBinary(buf))
	assert.Equal(t, hDecode, h)
}

func TestBlock_Encode_Decode(t *testing.T){
	b := &Block{
		Header: Header{
		Version:1,
		PrevBlock: types.RandomHash(),
		Timestamp: time.Now().UnixNano(),
		Height:32,
		Nonce:12450,
		},
		Transactions: nil,
	}

	buf := &bytes.Buffer{}
	assert.Nil(t, b.EncodeBinary(buf))

	bDecode := &Block{}
	assert.Nil(t, bDecode.DecodeBinary(buf))
	assert.Equal(t, bDecode, b)
}

func TestBlockHash(t *testing.T){
	b := &Block{
		Header: Header{
		Version:1,
		PrevBlock: types.RandomHash(),
		Timestamp: time.Now().UnixNano(),
		Height:32,
		Nonce:12450,
		},
		Transactions: []Transaction{},
	}

	h := b.Hash()
	assert.False(t, h.IsZero())

	fmt.Printf("%+v",h)
}