package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func newBlockChainWithGenesis(t *testing.T) *BlockChain {
	bc, err := NewBlockChain(randomBlock(0))

	assert.Nil(t, err)

	return bc
}

// Ideal scenario
func TestAddBlock(t *testing.T) {
	testingLen := 1000

	bc := newBlockChainWithGenesis(t)

	assert.Equal(t, bc.Height(), uint32(0))

	for i := range testingLen {
		// block := randomBlock(uint32(i + 1)) // No signature
		block := randomBlockWithSignature(t, uint32(i+1))
		assert.Nil(t, bc.AddBlock(block))
	}

	assert.Equal(t, bc.Height(), uint32(testingLen))
	assert.Equal(t, len(bc.headers), testingLen+1)
	assert.NotNil(t, bc.AddBlock(randomBlock(888))) // Validate block
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
}

// // uint32 overflow test
// func TestUint(t *testing.T) {
// 	var i uint32

// 	tx := make([]byte, 0)

// 	fmt.Println(i)
// 	fmt.Println(uint32(len(tx) - 1))
// }
