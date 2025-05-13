package core

import "io"

// On-chain data with this module
// TODO: Need signature
type Transaction struct {
	Data []byte // This can be any arbitrary data (Will be a VOTE DATA)
}

func (tx *Transaction) EncodeBinary(w io.Writer) error {
	return nil
}

func (tx *Transaction) DecodeBinary(r io.Reader) error {
	return nil
}
