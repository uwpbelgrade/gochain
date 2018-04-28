package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
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

// NewCoinbaseTransaction creates new coinbase transaction
func NewCoinbaseTransaction(to, data string, reward int) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward %s", to)
	}
	txin := TxInput{[]byte{}, -1, data}
	txout := TxOutput{reward, to}
	tx := &Transaction{[]byte{}, []TxInput{txin}, []TxOutput{txout}}
	tx.GenerateID()
	return tx
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

// IsCoinbase checks whether the transaction is coinbase
func (tx Transaction) IsCoinbase() bool {
	return len(tx.Vin) == 1 && len(tx.Vin[0].Txid) == 0 && tx.Vin[0].Vout == -1
}

// CanUnlockOutput checks if key can unlock output
func (txin *TxInput) CanUnlockOutput(key string) bool {
	return txin.ScriptSig == key
}

// CanOutputBeUnlocked checks if output can be unlocked with key
func (txout *TxOutput) CanOutputBeUnlocked(key string) bool {
	return txout.ScriptPubKey == key
}
