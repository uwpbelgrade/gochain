package core

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitChain(t *testing.T) {
	env := &EnvConfig{}
	os.Remove(env.GetDbFile())
	chain := InitChain(env, "address")
	assert.NotNil(t, chain)
	chain.db.Close()
}

func TestGetChain(t *testing.T) {
	env := &EnvConfig{}
	chain := GetChain(env)
	assert.NotNil(t, chain)
	chain.Log()
	chain.db.Close()
}

func TestIterateChain(t *testing.T) {
	env := &EnvConfig{}
	os.Remove(env.GetDbFile())
	chain := InitChain(env, "address")
	iterator := chain.Iterator()
	for {
		block := iterator.Next()
		assert.NotNil(t, block)
		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
	chain.db.Close()
}

func TestGetBalance(t *testing.T) {
	env := &EnvConfig{}
	os.Remove(env.GetDbFile())
	ws := NewWalletStore(env)
	wallet := ws.CreateWallet()
	address := string(wallet.GetAddress())
	chain := InitChain(env, address)
	balance := chain.GetBalance(address)
	assert.Equal(t, 50, balance)
	chain.db.Close()
}

func TestSendTransaction(t *testing.T) {
	env := &EnvConfig{}
	os.Remove(env.GetDbFile())
	ws := NewWalletStore(env)
	wallet1 := ws.CreateWallet()
	wallet2 := ws.CreateWallet()
	address1 := string(wallet1.GetAddress())
	address2 := string(wallet2.GetAddress())
	chain := InitChain(env, address1)
	assert.Equal(t, 50, chain.GetBalance(address1))
	chain.NewTransaction(address1, address2, 10)
	assert.Equal(t, 40, chain.GetBalance(address1))
	assert.Equal(t, 10, chain.GetBalance(address2))
	chain.db.Close()
}

func TestFailSendTransactionNotEnoughBalance(t *testing.T) {
	env := &EnvConfig{}
	os.Remove(env.GetDbFile())
	ws := NewWalletStore(env)
	wallet1 := ws.CreateWallet()
	wallet2 := ws.CreateWallet()
	address1 := string(wallet1.GetAddress())
	address2 := string(wallet2.GetAddress())
	chain := InitChain(env, address1)
	_, err := chain.NewTransaction(address1, address2, 60)
	assert.Equal(t, "not enough balance", err.Error())
	assert.Equal(t, 50, chain.GetBalance(address1))
	assert.Equal(t, 0, chain.GetBalance(address2))
	chain.db.Close()
}
