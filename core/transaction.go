package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
)

// TxInput transaction input
type TxInput struct {
	Txid      []byte
	Vout      int
	ScriptSig string
}

// TxOutput transaction output
type TxOutput struct {
	Value        int
	ScriptPubKey string
}

// Transaction model
type Transaction struct {
	ID   []byte
	Vin  []TxInput
	Vout []TxOutput
}

// GenerateID generates new id for transaction
func (tx *Transaction) GenerateID() {
	var encoded bytes.Buffer
	encoder := gob.NewEncoder(&encoded)
	err := encoder.Encode(tx)
	if err != nil {
		panic(err)
	}
	hash := sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}
