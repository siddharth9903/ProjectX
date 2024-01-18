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

	if(b.Height != v.bc.Height()+1 ){
		return fmt.Errorf("block is too high with height (%d) with hash (%s)", b.Height, b.Hash(BlockHasher{}))
	}

	prevHeader, err := v.bc.GetHeader(b.Height - 1);
	if err != nil { 
		return err
	}

	hash := BlockHasher{}.Hash(prevHeader);

	if b.PrevBlockHash != hash {
		return fmt.Errorf("the hash of prev block is invalid")
	}

	if err := b.Verify(); err != nil {
		return err
	}

	return nil
}