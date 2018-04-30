package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTxInLocking(t *testing.T) {
	ws := NewWalletStore(&EnvConfig{})
	wallet := ws.CreateWallet()
	address := string(wallet.GetAddress())
	txin := &TxInput{[]byte("1"), 0, nil, wallet.PublicKey}
	pubKeyHash, _ := PubKeyHash(address)
	assert.True(t, txin.CanUnlockOutput(pubKeyHash))
}

func TestTxOutLocking(t *testing.T) {
	ws := NewWalletStore(&EnvConfig{})
	wallet := ws.CreateWallet()
	address := string(wallet.GetAddress())
	txout := NewTxOutput(10, address)
	hash, _ := PubKeyHash(address)
	assert.True(t, txout.CanOutputBeUnlocked(hash))
}

func TestNewCoinbaseTransaction(t *testing.T) {
	ws := NewWalletStore(&EnvConfig{})
	w := ws.CreateWallet()
	tx := NewCoinbaseTransaction(string(w.GetAddress()), "data", 50)
	assert.Equal(t, 50, tx.Vout[0].Value)
	assert.Equal(t, true, tx.IsCoinbase())
	assert.Equal(t, "data", string(tx.Vin[0].PubKey))
}
