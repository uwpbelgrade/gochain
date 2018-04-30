package core

import (
	"bytes"
	"encoding/hex"
	"fmt"
)

// TxInput transaction input
type TxInput struct {
	Txid      []byte
	Vout      int
	Signature []byte
	PubKey    []byte
}

// CanUnlockOutput checks if key can unlock output
func (txin *TxInput) CanUnlockOutput(pubKeyHash []byte) bool {
	hash := RipeMd160Sha256(txin.PubKey)
	return bytes.Compare(pubKeyHash, hash) == 0
}

// Log logs txin
func (txin *TxInput) Log() {
	template := "\t\t[tx:%s][out:%d][sign:%x][pubkey:%x]\n"
	fmt.Printf(template, hex.EncodeToString(txin.Txid), txin.Vout, txin.Signature, txin.PubKey)
}
