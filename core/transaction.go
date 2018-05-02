package core

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
)

// Transaction model
type Transaction struct {
	ID   []byte
	Vin  []TxInput
	Vout []TxOutput
}

// NewCoinbaseTransaction creates new coinbase transaction
func NewCoinbaseTransaction(to, data string, reward int) *Transaction {
	txin := TxInput{[]byte{}, -1, nil, []byte(data)}
	txout := NewTxOutput(reward, to)
	tx := &Transaction{[]byte{}, []TxInput{txin}, []TxOutput{*txout}}
	tx.ID = tx.Hash()
	return tx
}

// Serialize serializes the transaction
func (tx *Transaction) Serialize() []byte {
	var encoded bytes.Buffer
	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	return encoded.Bytes()
}

// IsCoinbase checks whether the transaction is coinbase
func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Vin) == 1 && len(tx.Vin[0].Txid) == 0 && tx.Vin[0].Vout == -1
}

// Hash returns transaction hash
func (tx *Transaction) Hash() []byte {
	tcopy := *tx
	tcopy.ID = []byte{}
	hash := sha256.Sum256(tcopy.Serialize())
	return hash[:]
}

// AsSignaturePayload makes copy of transaction with required elements for signing
func (tx *Transaction) AsSignaturePayload() Transaction {
	var ins []TxInput
	var outs []TxOutput
	for _, i := range tx.Vin {
		ins = append(ins, TxInput{i.Txid, i.Vout, nil, nil})
	}
	for _, o := range tx.Vout {
		outs = append(outs, TxOutput{o.Value, o.PubKeyHash})
	}
	return Transaction{tx.ID, ins, outs}
}

// Sign signs transaction using private key
func (tx *Transaction) Sign(pk *ecdsa.PrivateKey, previousTxs map[string]Transaction) {
	if tx.IsCoinbase() {
		return
	}
	payload := tx.AsSignaturePayload()
	for ii, i := range payload.Vin {
		previousTx := previousTxs[hex.EncodeToString(i.Txid)]
		payload.Vin[ii].Signature = nil
		payload.Vin[ii].PubKey = previousTx.Vout[i.Vout].PubKeyHash
		payload.ID = payload.Hash()
		x, y, _ := ecdsa.Sign(rand.Reader, pk, payload.Hash())
		tx.Vin[ii].Signature = append(x.Bytes(), y.Bytes()...)
	}
}

// Verify verifies transaction
func (tx *Transaction) Verify(previousTxs map[string]Transaction) bool {
	curve := elliptic.P256()
	payload := tx.AsSignaturePayload()
	for ii, i := range tx.Vin {
		previousTx := previousTxs[hex.EncodeToString(i.Txid)]
		payload.Vin[ii].Signature = nil
		payload.Vin[ii].PubKey = previousTx.Vout[i.Vout].PubKeyHash
		payload.ID = payload.Hash()
		xSig := &big.Int{}
		ySig := &big.Int{}
		lenSig := len(i.Signature)
		xSig.SetBytes(i.Signature[:lenSig/2])
		xKey := &big.Int{}
		yKey := &big.Int{}
		lenKey := len(i.PubKey)
		ySig.SetBytes(i.Signature[lenSig/2:])
		xKey.SetBytes(i.PubKey[:lenKey/2])
		yKey.SetBytes(i.PubKey[lenKey/2:])
		rawPubKey := &ecdsa.PublicKey{Curve: curve, X: xKey, Y: yKey}
		if !ecdsa.Verify(rawPubKey, payload.ID, xSig, ySig) {
			return false
		}
	}
	return true
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
