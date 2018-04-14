package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPow(t *testing.T) {
	chain := InitChain()
	chain.AddBlock("data")
	block := chain.Get(chain.tip)
	assert.True(t, block.ValidatePOW(), "Hash is correct")
}
