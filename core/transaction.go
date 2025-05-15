package core

import (
	"fmt"

	"github.com/andantan/go-node/crypto"
)

// On-chain data with this module
type Transaction struct {
	Data []byte // This can be any arbitrary data (Will be a VOTE DATA)

	From      crypto.PublicKey
	Signature *crypto.Signature
}

func (tx *Transaction) Sign(privKey crypto.PrivateKey) error {
	sig, err := privKey.Sign(tx.Data)

	if err != nil {
		return err
	}

	tx.From = privKey.PublicKey()
	tx.Signature = sig

	return nil
}

func (tx *Transaction) Verify() error {
	if tx.Signature == nil {
		return fmt.Errorf("transaction has no signature")
	}

	if !tx.Signature.Verify(tx.From, tx.Data) {
		return fmt.Errorf("invaild transaction signature")
	}

	return nil
}
