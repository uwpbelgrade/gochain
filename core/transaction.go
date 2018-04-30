package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
)

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
	txin := TxInput{[]byte{}, -1, nil, []byte(data)}
	txout := NewTxOutput(reward, to)
	tx := &Transaction{[]byte{}, []TxInput{txin}, []TxOutput{*txout}}
	tx.GenerateID()
	return tx
}

// GenerateID generates new id for transaction
func (tx *Transaction) GenerateID() {
	var encoded bytes.Buffer
	encoder := gob.NewEncoder(&encoded)
	encoder.Encode(tx)
	hash := sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}

// IsCoinbase checks whether the transaction is coinbase
func (tx Transaction) IsCoinbase() bool {
	return len(tx.Vin) == 1 && len(tx.Vin[0].Txid) == 0 && tx.Vin[0].Vout == -1
}

// Log logs transaction
func (tx *Transaction) Log() {
	template := "[txid:%s]\n"
	fmt.Printf(template, hex.EncodeToString(tx.ID))
	fmt.Printf("\tinputs:\n")
	for _, i := range tx.Vin {
		i.Log()
	}
	fmt.Printf("\toutputs:\n")
	for _, o := range tx.Vout {
		o.Log()
	}
}
