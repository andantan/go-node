package core

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/andantan/go-node/crypto"
	"github.com/andantan/go-node/types"
)

type Header struct {
	Version       uint32 // Decode base
	DataHash      types.Hash
	PrevBlockHash types.Hash
	TimeStamp     int64
	// Will add space field for appendable something
	Height uint32
	// Nonce  uint64	// Not need Nonce
}

func (h *Header) Bytes() []byte {
	buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)

	enc.Encode(h)

	return buf.Bytes()
}

type Block struct {
	*Header
	Transactions []Transaction
	// Need validator which as PublicKey the reason for consensus
	// elect leader and leader as privilege to propose network
	Validator crypto.PublicKey
	Signature *crypto.Signature

	hash types.Hash // Cached version of the header heash
}

func NewBlock(h *Header, txx []Transaction) *Block {
	return &Block{
		Header:       h,
		Transactions: txx,
	}
}

func (b *Block) AddTransaction(tx *Transaction) {
	b.Transactions = append(b.Transactions, *tx)
}

func (b *Block) Sign(privKey crypto.PrivateKey) error {
	sig, err := privKey.Sign(b.Header.Bytes())

	if err != nil {
		return err
	}

	b.Validator = privKey.PublicKey()
	b.Signature = sig

	return nil
}

// Verify Signature
func (b *Block) Verify() error {
	if b.Signature == nil {
		return fmt.Errorf("block has no signature")
	}

	if !b.Signature.Verify(b.Validator, b.Header.Bytes()) {
		return fmt.Errorf("block has invalid signature")
	}

	for _, tx := range b.Transactions {
		if err := tx.Verify(); err != nil {
			return err
		}
	}

	return nil
}

func (b *Block) Decode(dec Decoder[*Block]) error {
	return dec.Decode(b)
}

func (b *Block) Encode(enc Encoder[*Block]) error {
	return enc.Encode(b)
}

// Capsulate hash function
func (b *Block) Hash(hasher Hasher[*Header]) types.Hash {
	if b.hash.IsZero() {
		b.hash = hasher.Hash(b.Header)
	}

	return b.hash
}
