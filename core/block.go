package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"time"

	"github.com/andantan/go-node/crypto"
	"github.com/andantan/go-node/types"
)

type Header struct {
	Version       uint32     // Decode base
	DataHash      types.Hash // Hashs of All of transaction
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
	Transactions []*Transaction
	// Need validator which as PublicKey the reason for consensus
	// elect leader and leader as privilege to propose network
	Validator crypto.PublicKey
	Signature *crypto.Signature

	hash types.Hash // Cached version of the header heash
}

func NewBlock(h *Header, txx []*Transaction) (*Block, error) {
	return &Block{
		Header:       h,
		Transactions: txx,
	}, nil
}

func NewBlockFromPrevHeader(prevHeader *Header, txx []*Transaction) (*Block, error) {
	dataHash, err := CalculateDataHash(txx)

	if err != nil {
		return nil, err
	}

	header := &Header{
		Version:       1,
		DataHash:      dataHash,
		PrevBlockHash: BlockHasher{}.Hash(prevHeader),
		TimeStamp:     time.Now().UnixNano(),
		Height:        prevHeader.Height + 1,
	}

	return NewBlock(header, txx)

}

func (b *Block) AddTransaction(tx *Transaction) {
	b.Transactions = append(b.Transactions, tx)
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

	dataHash, err := CalculateDataHash(b.Transactions)

	if err != nil {
		return err
	}

	if dataHash != b.DataHash {
		return fmt.Errorf("block (%s) has an invalid data hash",
			b.Hash(BlockHasher{}))
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

func CalculateDataHash(txx []*Transaction) (hash types.Hash, err error) {

	buf := &bytes.Buffer{}

	for _, tx := range txx {
		if err = tx.Encode(NewGobTxEncoder(buf)); err != nil {
			return
		}
	}

	hash = sha256.Sum256(buf.Bytes())

	return
}
