package core

import (
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
}
