package core

import (
	"bytes"
	"encoding/gob"
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
		hashes = append(hashes, transaction.Serialize())
	}
	if len(hashes) == 0 {
		return []byte{}
	}
	mtree := NewMerkleTree(hashes)
	return mtree.Root.Data
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
	template := "BLOCK >>>> \nPrevious hash: %x \nData: %x " +
		"\nTimestamp: %d [%s] \nNonce: %d \nTransactions:\n"
	fmt.Printf(template, block.PrevBlockHash, block.Hash,
		block.Timestamp, time.Unix(block.Timestamp, 0), block.Nonce)
	for _, t := range block.Transactions {
		t.Log()
	}
	fmt.Printf("\n")
}
