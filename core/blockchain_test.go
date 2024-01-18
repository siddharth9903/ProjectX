package core

import (
	"ProjectX/types"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newBlockchainWithGenesis(t *testing.T) *Blockchain {
	bc, err := NewBlockchain(randomBlock(0, types.Hash{}))

	assert.Nil(t,err)
	return bc
}

func TestNewBlockchain(t *testing.T) {
	bc := newBlockchainWithGenesis(t)

	assert.NotNil(t,bc)
}

func getPreviousBlockHash(t *testing.T, bc *Blockchain, height uint32) types.Hash {
	prevHeader,err := bc.GetHeader(height - 1);

	assert.Nil(t, err);
	prevBlockHash := BlockHasher{}.Hash(prevHeader);
	// assert.NotNil(t, prevBlockHash)
	return prevBlockHash
}
func TestAddBlock(t *testing.T) {
	bc := newBlockchainWithGenesis(t);
	header, err := bc.GetHeader(bc.Height());
	fmt.Println("header",header)
	assert.Nil(t, err);
	assert.NotNil(t, header);

	lenBlocks := 1000;
	for i := 0; i < lenBlocks; i++ {
		prevHash := getPreviousBlockHash(t, bc, uint32(i+1));

		b := randomBlockWithSignature(uint32(i+1),prevHash)
		err := bc.AddBlock(b);
		assert.Nil(t, err);

		header, err := bc.GetHeader(bc.Height());
		assert.Nil(t, err);
		assert.NotNil(t, header);
	}
}

func TestHasBlock(t *testing.T){
	bc := newBlockchainWithGenesis(t)
	assert.True(t, bc.HasBlock(0))
	assert.False(t, bc.HasBlock(1))
}