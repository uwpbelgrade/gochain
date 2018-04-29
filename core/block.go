package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"time"
)

// Block holds transactions in chain
type Block struct {
	Timestamp     int64
	Transactions  []*Transaction
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int
}

// NewBlock creates new block
func NewBlock(transactions []*Transaction, prevBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(), transactions, prevBlockHash, []byte{}, 0}
	block.POW()
	return block
}

// HashTransactions makes hash of all transaction ids
func (block *Block) HashTransactions() []byte {
	var hashes [][]byte
	for _, transaction := range block.Transactions {
		hashes = append(hashes, transaction.ID)
	}
	hash := sha256.Sum256(bytes.Join(hashes, []byte{}))
	return hash[:]
}

// Serialize serializes block using encoder
func (block *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(block)
	if err != nil {
		panic(err)
	}
	return result.Bytes()
}

// Deserialize deserializes bytes to block
func Deserialize(data []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&block)
	if err != nil {
		panic(err)
	}
	return &block
}

// Log prints block info
func (block *Block) Log() {
	template := `
	BLOCK >>>>
	Previous hash: %x
	Data: %x
	Timestamp: %d [%s]
	Nonce: %d
	Transactions:
	`
	fmt.Printf(template, block.PrevBlockHash, block.Hash,
		block.Timestamp, time.Unix(block.Timestamp, 0), block.Nonce)

	for _, t := range block.Transactions {
		fmt.Printf("t[ID]: %s \n", hex.EncodeToString(t.ID))
	}
	fmt.Printf("\n")
}
