package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewWallet(t *testing.T) {
	wallet := NewWallet()
	assert.NotNil(t, wallet)
	assert.NotNil(t, wallet.PrivateKey)
	assert.NotNil(t, wallet.PublicKey)
	assert.True(t, wallet.IsValidAddress())
}

func TestGetAddress(t *testing.T) {
	wallet := NewWallet()
	assert.NotNil(t, wallet.GetAddress())
}
