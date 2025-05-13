package core

import (
	"bytes"
	"testing"
	"time"

	"github.com/andantan/go-node/types"
	"github.com/stretchr/testify/assert"
)

func TestHeader_Encode_Decode(t *testing.T) {
	h := &Header{
		Version:   1,
		PrevBlock: types.RandomHash(), // Hashed random PrevBlock
		TimeStamp: uint64(time.Now().UnixNano()),
		Height:    10, // Tree blocks: test genesis block
		Nonce:     989394,
	}

	buf := &bytes.Buffer{}

	assert.Nil(t, h.EncodeBinary(buf))

	hDecode := &Header{}

	assert.Nil(t, hDecode.DecodeBinary(buf))
	assert.Equal(t, h, hDecode)
}

func TestBlock_Encode_Decode(t *testing.T) {
	b := &Block{
		Header: Header{
			Version:   1,
			PrevBlock: types.RandomHash(),
			TimeStamp: uint64(time.Now().UnixNano()),
			Height:    10,
			Nonce:     989394,
		},
		Transactions: nil,
	}

	buf := &bytes.Buffer{}

	assert.Nil(t, b.EncodeBinary(buf))

	bDecode := &Block{}

	assert.Nil(t, bDecode.DecodeBinary(buf))
	assert.Equal(t, b, bDecode)

	// fmt.Printf("%+v\n", bDecode)
}

func TestBlockHash(t *testing.T) {
	b := &Block{
		Header: Header{
			Version:   1,
			PrevBlock: types.RandomHash(),
			TimeStamp: uint64(time.Now().UnixNano()),
			Height:    10,
			Nonce:     989394,
		},
		Transactions: []Transaction{},
	}

	h := b.Hash()

	// fmt.Println(h)
	assert.False(t, h.IsZero())

}
