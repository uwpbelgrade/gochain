package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCoinbaseTransaction(t *testing.T) {
	tx := NewCoinbaseTransaction("abcd", "data", 50)
	assert.Equal(t, 50, tx.Vout[0].Value)
	assert.Equal(t, true, tx.IsCoinbase())
	assert.Equal(t, "datsa", tx.Vin[0].ScriptSig)
}
