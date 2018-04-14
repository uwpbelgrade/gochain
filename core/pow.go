package core

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"strings"
)

// Difficulty of challenge
const Difficulty = 4

// POW finds POW: sha256 hash that starts with N (difficulty) zeros
func (block *Block) POW() {
	prefixTarget := strings.Repeat("0", Difficulty)
	for nonce := 0; nonce < math.MaxInt64; nonce++ {
		hash := sha256.Sum256(block.join(nonce))
		prefix := fmt.Sprintf("%x", hash)[0:Difficulty]
		if prefixTarget == prefix {
			fmt.Printf("pow found, hash %x, nonce %d \n", hash, nonce)
			block.Nonce = nonce
			block.Hash = hash[:]
			return
		}
	}
	panic(nil)
}

// ValidatePOW block hash
func (block *Block) ValidatePOW() bool {
	actual := fmt.Sprintf("%x", block.Hash)[0:Difficulty]
	required := strings.Repeat("0", Difficulty)
	return strings.HasPrefix(actual, required)
}

func (block *Block) join(nonce int) []byte {
	return bytes.Join([][]byte{
		block.PrevBlockHash,
		block.Data,
		block.Hash,
		[]byte(fmt.Sprintf("%x", nonce)),
	}, []byte{})
}
