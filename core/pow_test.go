package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPow(t *testing.T) {
	chain := InitChain("address")
	chain.AddBlock([]*Transaction{DemoTransaction()})
	block := chain.Get(chain.tip)
	assert.True(t, block.ValidatePOW(), "Hash is correct")
}
