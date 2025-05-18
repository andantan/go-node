package core

import (
	"testing"
	"time"

	"github.com/andantan/go-node/crypto"
	"github.com/andantan/go-node/types"
	"github.com/stretchr/testify/assert"
)

func TestSignBlock(t *testing.T) {
	privKey := crypto.GeneratePrivateKey()
	b := randomBlock(t, 0, types.Hash{}) // Genesis

	assert.Nil(t, b.Sign(privKey))
	assert.NotNil(t, b.Signature)
}

func TestVerifyBlock(t *testing.T) {
	privKey := crypto.GeneratePrivateKey()
	b := randomBlock(t, 0, types.Hash{}) // Genesis

	assert.Nil(t, b.Sign(privKey))
	assert.Nil(t, b.Verify())

	otherPrivKey := crypto.GeneratePrivateKey()

	// Tempering public-key
	b.Validator = otherPrivKey.PublicKey()

	assert.NotNil(t, b.Verify())

	b.Height = 100

	assert.NotNil(t, b.Verify())
}

func randomBlock(t *testing.T, height uint32, PrevBlockHash types.Hash) *Block {
	privKey := crypto.GeneratePrivateKey()
	tx := randomTxWithSignature(t)
	header := &Header{
		Version:       1,
		PrevBlockHash: PrevBlockHash,
		Height:        height,
		TimeStamp:     time.Now().UnixNano(),
	}

	b, err := NewBlock(header, []Transaction{tx})

	assert.Nil(t, err)

	dataHash, err := CalculateDataHash(b.Transactions)

	assert.Nil(t, err)

	b.Header.DataHash = dataHash

	assert.Nil(t, b.Sign(privKey))

	return b
}
