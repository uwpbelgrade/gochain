package core

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPow(t *testing.T) {
	chain := InitChain()
	chain.AddBlock("data")
	_, hash := Work(chain.Blocks[len(chain.Blocks)-1])
	assert.Equal(t, fmt.Sprintf("%x", hash)[0:Difficulty], strings.Repeat("0", Difficulty), "Hash is correct")
}
