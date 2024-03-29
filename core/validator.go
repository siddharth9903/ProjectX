package core

import "fmt"

type Validator interface {
	ValidateBlock(b *Block) error
}

type BlockValidator struct {
	bc *Blockchain
}

func NewBlockValidator(bc *Blockchain) *BlockValidator {
	return &BlockValidator{
		bc: bc,
	}
}

func (v *BlockValidator) ValidateBlock(b *Block) error {

	if v.bc.HasBlock(b.Height){
		return fmt.Errorf("chain already contains block (%d) with hash (%s)", b.Height, b.Hash(BlockHasher{}))
	}

	if b.Height != v.bc.Height()+1 {
		return fmt.Errorf("block (%s) with height (%d) is too high => current height (%d)", b.Hash(BlockHasher{}), b.Height, v.bc.Height())
	}

	prevHeader, err := v.bc.GetHeader(b.Height - 1);
	if err != nil { 
		return err
	}

	hash := BlockHasher{}.Hash(prevHeader);

	if b.PrevBlockHash != hash {
		return fmt.Errorf("the hash of prev block is invalid")
	}

	dataHash,err := CalculateDataHash(b.Transactions)
	if err != nil {
		return err
	}

	if dataHash != b.DataHash {
		return fmt.Errorf("block %s has invalid datahash", b.Hash(BlockHasher{}))
	}

	if err := b.Verify(); err != nil {
		return err
	}

	return nil
}