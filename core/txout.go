package core

import (
	"bytes"
	"fmt"

	"github.com/mr-tron/base58/base58"
)

// TxOutput transaction output
type TxOutput struct {
	Value      int
	PubKeyHash []byte
}

// NewTxOutput creates new TxOutput and locks it for address
func NewTxOutput(value int, address string) *TxOutput {
	txo := &TxOutput{value, nil}
	txo.LockOutput(address)
	return txo
}

// LockOutput locks output with public key hash of given address
func (txout *TxOutput) LockOutput(address string) error {
	pubKeyHash, _ := base58.Decode(address)
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-AddressChecksumLength]
	txout.PubKeyHash = pubKeyHash
	return nil
}

// CanOutputBeUnlocked checks if output can be unlocked with key
func (txout *TxOutput) CanOutputBeUnlocked(pubKeyHash []byte) bool {
	return bytes.Compare(txout.PubKeyHash, pubKeyHash) == 0
}

// Log logs txout
func (txout *TxOutput) Log() {
	template := "\t\t[val:%d][pubkeyhash:%x]\n"
	fmt.Printf(template, txout.Value, txout.PubKeyHash)
}
