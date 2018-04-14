package core

import (
	"bytes"
	"encoding/gob"
	"time"
)

// Block holds transactions in chain
type Block struct {
	Timestamp     int64
	Data          []byte
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int
}

// NewBlock creates new block in chain
func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(), []byte(data), prevBlockHash, []byte{}, 0}
	Work(block)
	return block
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
