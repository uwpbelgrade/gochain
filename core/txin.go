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
	// fmt.Printf("PUBKEY LEN: %d", len(txin.PubKey))
	if len(txin.PubKey) != 64 {
		template := "\t\t[tx:%s]\n\t\t[out:%d] [sign:%d] [pubkey:%s]\n"
		fmt.Printf(template, hex.EncodeToString(txin.Txid), txin.Vout, len(txin.Signature), string(txin.PubKey))
	} else {
		template := "\t\t[tx:%s]\n\t\t[out:%d] [sign:%d] [address:%s]\n"
		address := string(GetAddressFromPublicKey(txin.PubKey))
		fmt.Printf(template, hex.EncodeToString(txin.Txid), txin.Vout, len(txin.Signature), address)
	}
}
