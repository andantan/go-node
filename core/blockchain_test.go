package core

import (
	"testing"

	"github.com/andantan/go-node/types"
	"github.com/stretchr/testify/assert"
)

// Ideal scenario
func TestAddBlock(t *testing.T) {
	testingLen := 1000

	bc := newBlockChainWithGenesis(t)

	assert.Equal(t, bc.Height(), uint32(0))

	for i := range testingLen {
		// block := randomBlock(uint32(i + 1)) // No signature
		block := randomBlockWithSignature(t,
			uint32(i+1), getPrevBlockHash(t, bc, uint32(i+1)))
		assert.Nil(t, bc.AddBlock(block))
	}

	assert.Equal(t, bc.Height(), uint32(testingLen))
	assert.Equal(t, len(bc.headers), testingLen+1)
	assert.NotNil(t, bc.AddBlock(randomBlock(89, types.Hash{}))) // Validate block
}

func TestNewBlockChain(t *testing.T) {
	bc := newBlockChainWithGenesis(t)

	assert.NotNil(t, bc.validator)
	assert.Equal(t, bc.Height(), uint32(0)) // GenesisBlock height must be 0

	// fmt.Println(bc.Height())
}

func TestHasBlock(t *testing.T) {
	bc := newBlockChainWithGenesis(t)

	assert.True(t, bc.HasBlock(0))
	assert.False(t, bc.HasBlock(44))
}

func TestGetHeader(t *testing.T) {
	testingLen := 100

	bc := newBlockChainWithGenesis(t)

	for i := range testingLen {
		blockHeight := uint32(i + 1)

		block := randomBlockWithSignature(t, blockHeight, getPrevBlockHash(t, bc, uint32(i+1)))
		assert.Nil(t, bc.AddBlock(block))

		header, err := bc.GetHeader(blockHeight)
		assert.Nil(t, err)
		assert.Equal(t, header, block.Header)
	}
}

func TestAddBlockToHigh(t *testing.T) {
	bc := newBlockChainWithGenesis(t)

	assert.Nil(t, bc.AddBlock(randomBlockWithSignature(t, 1, getPrevBlockHash(t, bc, uint32(1)))))
	assert.NotNil(t, bc.AddBlock(randomBlockWithSignature(t, 3, types.Hash{})))
}

func newBlockChainWithGenesis(t *testing.T) *BlockChain {
	bc, err := NewBlockChain(randomBlock(0, types.Hash{}))

	assert.Nil(t, err)

	return bc
}

// Find previous block hash based on Height
func getPrevBlockHash(t *testing.T, bc *BlockChain, height uint32) types.Hash {
	prevHeader, err := bc.GetHeader(height - 1)

	assert.Nil(t, err)

	return BlockHasher{}.Hash(prevHeader)
}

// // uint32 overflow test
// func TestUint(t *testing.T) {
// 	var i uint32

// 	tx := make([]byte, 0)

// 	fmt.Println(i)
// 	fmt.Println(uint32(len(tx) - 1))
// }
