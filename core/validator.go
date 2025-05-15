package core

import "fmt"

type Validator interface {
	ValidateBlock(*Block) error
}

type BlockValidator struct {
	bc *BlockChain
}

func NewBlockValidator(bc *BlockChain) *BlockValidator {
	return &BlockValidator{
		bc: bc,
	}
}

func (v *BlockValidator) ValidateBlock(b *Block) error {
	if v.bc.HasBlock(b.Height) {
		return fmt.Errorf("chain already contains block (%d) with hash (%s)",
			b.Height, b.Hash(BlockHasher{}))
	}

	// Block chaining
	if b.Height != v.bc.Height()+1 {
		return fmt.Errorf("block (%s) too high can't chaining",
			b.Hash(BlockHasher{}))
	}

	prevHeader, err := v.bc.GetHeader(b.Height - 1)

	if err != nil {
		return err
	}

	hash := BlockHasher{}.Hash(prevHeader)

	if hash != b.PrevBlockHash {
		return fmt.Errorf("the hash of the previous block (%s) is invalid", b.PrevBlockHash)
	}

	if err := b.Verify(); err != nil {
		return err
	}

	return nil
}
