package core

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type data struct {
	ws       *WalletStore
	wallet   *Wallet
	wallet2  *Wallet
	address1 string
	address2 string
	chain    *Blockchain
}

func newData(removeFile bool) *data {
	var w1, w2 *Wallet
	var a1, a2 string
	env := &EnvConfig{}
	wstore := NewWalletStore(env)
	if removeFile {
		os.Remove(env.GetDbFile())
	}
	w1 = wstore.CreateWallet()
	a1 = string(w1.GetAddress())
	w2 = wstore.CreateWallet()
	a2 = string(w2.GetAddress())
	chain := InitChain(env, a1)
	return &data{wstore, w1, w2, a1, a2, chain}
}

func TestInitChain(t *testing.T) {
	d := newData(true)
	assert.NotNil(t, d.chain)
	d.chain.db.Close()
}

func TestGetChain(t *testing.T) {
	env := &EnvConfig{}
	chain := GetChain(env)
	assert.NotNil(t, chain)
	chain.Log()
	chain.db.Close()
}

func TestIterateChain(t *testing.T) {
	d := newData(true)
	iterator := d.chain.Iterator()
	for {
		block := iterator.Next()
		assert.NotNil(t, block)
		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
	d.chain.db.Close()
}

func TestGetBalance(t *testing.T) {
	d := newData(true)
	balance := d.chain.GetBalance(d.address1)
	assert.Equal(t, 50, balance)
	d.chain.db.Close()
}

// func TestSendTransaction(t *testing.T) {
// 	d := newData(true)
// 	assert.Equal(t, 50, d.chain.GetBalance(d.address1))
// 	d.chain.Send(d.wallet, d.address2, 10)
// 	// tx, _ := d.chain.NewTransaction(d.address1, d.address2, 10)
// 	// d.chain.SignTransaction(&d.wallet.PrivateKey, tx)
// 	// d.chain.AddBlock([]*Transaction{tx})
// 	assert.Equal(t, 40, d.chain.GetBalance(d.address1))
// 	assert.Equal(t, 10, d.chain.GetBalance(d.address2))
// 	d.chain.db.Close()
// }

func TestFailSendTransactionNotEnoughBalance(t *testing.T) {
	d := newData(true)
	_, err := d.chain.NewTransaction(d.address1, d.address2, 60)
	assert.Equal(t, "not enough balance", err.Error())
	assert.Equal(t, 50, d.chain.GetBalance(d.address1))
	assert.Equal(t, 0, d.chain.GetBalance(d.address2))
	d.chain.db.Close()
}

func TestFailSendNotSigned(t *testing.T) {
	d := newData(true)
	tx, _ := d.chain.NewTransaction(d.address1, d.address2, 1)
	_, err := d.chain.AddBlock([]*Transaction{tx})
	assert.Contains(t, err.Error(), "invalid transaction")
	d.chain.db.Close()
}

func TestGetTransaction(t *testing.T) {
	d := newData(true)
	tx, _ := d.chain.NewTransaction(d.address1, d.address2, 1)
	d.chain.SignTransaction(&d.wallet.PrivateKey, tx)
	d.chain.AddBlock([]*Transaction{tx})
	txDb, _ := d.chain.GetTransaction(tx.ID)
	assert.Equal(t, tx, &txDb)
	d.chain.db.Close()
}

func TestFailGetTransacationUnknown(t *testing.T) {
	d := newData(true)
	_, err := d.chain.GetTransaction([]byte{0x00})
	assert.Contains(t, err.Error(), "transaction not found")
	d.chain.db.Close()
}
