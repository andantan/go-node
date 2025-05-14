package core

type BlockChain struct {
	store     Storage
	headers   []*Header
	validator Validator // Validator depends on BlockChain
}

func NewBlockChain(genesis *Block) (*BlockChain, error) {
	bc := &BlockChain{
		headers: []*Header{},
		store:   NewMemoryStore(),
	}

	bc.validator = NewBlockValidator(bc)
	err := bc.addBlockWithoutValidation(genesis)

	return bc, err
}

func (bc *BlockChain) SetValidator(v Validator) {
	bc.validator = v
}

// Validate a Block
func (bc *BlockChain) AddBlock(b *Block) error {
	if err := bc.validator.ValidateBlock(b); err != nil {
		return err
	}

	// Validation already done
	return bc.addBlockWithoutValidation(b)
}

func (bc *BlockChain) HasBlock(height uint32) bool {
	return height <= bc.Height()
}

// example: GenesisBlock [0, 1, 2, 3] => 4 len
// example: GenesisBlock [0, 1, 2, 3] => 3 height
func (bc *BlockChain) Height() uint32 {
	return uint32(len(bc.headers) - 1) // TestUint debugging overflows runtime error
}

// Internal addBlock for genesis block
func (bc *BlockChain) addBlockWithoutValidation(b *Block) error {
	bc.headers = append(bc.headers, b.Header)

	return bc.store.Put(b)
}
