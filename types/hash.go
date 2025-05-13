package types

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

// Tx | Header | ... -> SHA-256 -> 32Byte
type Hash [32]uint8

func (h Hash) IsZero() bool {
	for i := range 32 {
		if h[i] != 0 {
			return false
		}
	}

	return true
}

func (h Hash) ToSlice() []byte {
	b := make([]byte, 32)

	for i := range 32 {
		b[i] = h[i]
	}

	return b
}

func (h Hash) String() string {
	return hex.EncodeToString(h.ToSlice())
}

func HashFromBytes(b []byte) Hash {
	if len(b) != 32 {
		msg := fmt.Sprintf("Given bytes with length %d should be 32", len(b))

		// System can not continue
		panic(msg)
	}

	var value [32]uint8

	// Byte slice element by element
	// Not clone
	for i := range 32 {
		value[i] = b[i]
	}

	return Hash(value)
}

func RandomBytes(size int) []byte {
	token := make([]byte, size)

	// Fill cryptographically secure random bytes
	rand.Read(token)

	return token
}

func RandomHash() Hash {
	return HashFromBytes(RandomBytes(32))
}
