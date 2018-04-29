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
	chain := InitChain(env, "address")
	balance := chain.GetBalance("address")
	assert.Equal(t, 50, balance)
	chain.db.Close()
}

func TestSendTransaction(t *testing.T) {
	env := &EnvConfig{}
	os.Remove(env.GetDbFile())
	chain := InitChain(env, "a")
	chain.NewTransaction("a", "b", 10)
	assert.Equal(t, 40, chain.GetBalance("a"))
	assert.Equal(t, 10, chain.GetBalance("b"))
	chain.db.Close()
}

func TestFailSendTransactionNotEnoughBalance(t *testing.T) {
	env := &EnvConfig{}
	os.Remove(env.GetDbFile())
	chain := InitChain(env, "a")
	_, err := chain.NewTransaction("a", "b", 60)
	assert.Equal(t, "not enough balance", err.Error())
	assert.Equal(t, 50, chain.GetBalance("a"))
	assert.Equal(t, 0, chain.GetBalance("b"))
	chain.db.Close()
}
