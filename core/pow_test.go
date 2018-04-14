package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPow(t *testing.T) {
	chain := InitChain()
	chain.AddBlock("data")
	block := chain.Blocks[len(chain.Blocks)-1]
	Work(block)
	assert.True(t, Validate(block), "Hash is correct")
}
