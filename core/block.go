package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"io"

	"github.com/andantan/go-node/types"
)

type Header struct {
	Version   uint32 // Decode base
	PrevBlock types.Hash
	TimeStamp uint64
	// Will add space field for appendable something
	Height uint32
	Nonce  uint64
}

// Encode all fields to byte slice
// Extendable: Not return byte slice
// Order: LittleEndian -> LSB first
// 0x12345678 -> [0x78, 0x56, 0x34, 0x12] 일단 보류
func (h *Header) EncodeBinary(w io.Writer) error {
	if err := binary.Write(w, binary.LittleEndian, &h.Version); err != nil {
		return err
	}

	if err := binary.Write(w, binary.LittleEndian, &h.PrevBlock); err != nil {
		return err
	}

	if err := binary.Write(w, binary.LittleEndian, &h.TimeStamp); err != nil {
		return err
	}

	if err := binary.Write(w, binary.LittleEndian, &h.Height); err != nil {
		return err
	}

	return binary.Write(w, binary.LittleEndian, &h.Nonce)
}

// If system get byte slice then decode to Header
func (h *Header) DecodeBinary(r io.Reader) error {
	if err := binary.Read(r, binary.LittleEndian, &h.Version); err != nil {
		return err
	}

	if err := binary.Read(r, binary.LittleEndian, &h.PrevBlock); err != nil {
		return err
	}

	if err := binary.Read(r, binary.LittleEndian, &h.TimeStamp); err != nil {
		return err
	}

	if err := binary.Read(r, binary.LittleEndian, &h.Height); err != nil {
		return err
	}

	return binary.Read(r, binary.LittleEndian, &h.Nonce)
}

type Block struct {
	Header
	Transactions []Transaction

	hash types.Hash // Cached version of the header heash
}

// Each time calling hash
// if hash is empty then execute Block.Hash()
// Set the hash value and return value
func (b *Block) Hash() types.Hash {
	buf := &bytes.Buffer{}

	b.Header.EncodeBinary(buf)

	if b.hash.IsZero() {
		b.hash = types.Hash(sha256.Sum256(buf.Bytes()))
	}

	return b.hash
}

// Cannot encode if the binary does not hashed
func (b *Block) EncodeBinary(w io.Writer) error {
	if err := b.Header.EncodeBinary(w); err != nil {
		return err
	}

	for _, tx := range b.Transactions {
		if err := tx.EncodeBinary(w); err != nil {
			return err
		}
	}

	return nil
}

// Cannot decode if the binary does not hashed
func (b *Block) DecodeBinary(r io.Reader) error {
	if err := b.Header.DecodeBinary(r); err != nil {
		return err
	}

	for _, tx := range b.Transactions {
		if err := tx.DecodeBinary(r); err != nil {
			return err
		}
	}

	return nil
}
