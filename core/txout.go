package core

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"

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

// Serialize serializes TxOutputs
func (txout *TxOutput) Serialize() []byte {
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	err := enc.Encode(txout)
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}

// Deserialize deserializes bytes to TxOutputs
func (txout *TxOutput) Deserialize(data []byte) {
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(txout)
	if err != nil {
		panic(err)
	}
}

// SerializeOutputs deserializes outputs
func SerializeOutputs(data []TxOutput) []byte {
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	err := enc.Encode(data)
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}

// DeserializeOutputs deserializes outputs
func DeserializeOutputs(data []byte) []TxOutput {
	var outputs []TxOutput
	dec := gob.NewDecoder(bytes.NewReader(data))
	err := dec.Decode(&outputs)
	if err != nil {
		log.Panic(err)
	}
	return outputs
}

// Log logs txout
func (txout *TxOutput) Log() {
	template := "\t\t[val:%d] [address:%s]\n"
	address := string(GetAddressFromPublicKeyHash(txout.PubKeyHash))
	fmt.Printf(template, txout.Value, address)
}
