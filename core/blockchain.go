package core

import (
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"
)

// TODO(@andantan): Change BlockChain struct to interface
type BlockChain struct {
	store     Storage
	lock      sync.RWMutex
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

func (bc *BlockChain) GetHeader(height uint32) (*Header, error) {
	if height > bc.Height() {
		return nil, fmt.Errorf("given height (%d) too high", height)
	}
	bc.lock.Lock()
	defer bc.lock.Unlock()

	return bc.headers[int(height)], nil
}

func (bc *BlockChain) HasBlock(height uint32) bool {
	return height <= bc.Height()
}

// example: GenesisBlock [0, 1, 2, 3] => 4 len
// example: GenesisBlock [0, 1, 2, 3] => 3 height
func (bc *BlockChain) Height() uint32 {
	bc.lock.RLock()
	defer bc.lock.RUnlock()
	return uint32(len(bc.headers) - 1) // TestUint debugging overflows runtime error
}

// Internal addBlock for genesis block
func (bc *BlockChain) addBlockWithoutValidation(b *Block) error {
	bc.lock.Lock()
	bc.headers = append(bc.headers, b.Header)
	bc.lock.Unlock()

	logrus.WithFields(logrus.Fields{
		"height": b.Height,
		"hash":   b.Hash(BlockHasher{}),
	}).Info("adding new block")

	return bc.store.Put(b)
}
