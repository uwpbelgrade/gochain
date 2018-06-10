package core

import (
	"os"
	"testing"

	"github.com/boltdb/bolt"
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
	var nodeID = "1"
	var w1, w2 *Wallet
	var a1, a2 string
	env := &EnvConfig{}
	wstore := NewWalletStore(env, nodeID)
	if removeFile {
		os.Remove(env.GetDbFile(nodeID))
	}
	w1 = wstore.CreateWallet()
	a1 = string(w1.GetAddress())
	w2 = wstore.CreateWallet()
	a2 = string(w2.GetAddress())
	chain := InitChain(env, a1, nodeID)
	return &data{wstore, w1, w2, a1, a2, chain}
}

func TestInitChain(t *testing.T) {
	d := newData(true)
	assert.NotNil(t, d.chain)
	d.chain.db.Close()
}

func TestGetBestHeight(t *testing.T) {
	var nodeID = "1"
	env := &EnvConfig{}
	db, err := bolt.Open(env.GetDbFile(nodeID), 0600, nil)
	if err != nil {
		panic(err)
	}
	bestHeight := GetBestHeight(db, env)
	assert.NotNil(t, bestHeight)
	db.Close()
}

func TestGetChain(t *testing.T) {
	var nodeID = "1"
	env := &EnvConfig{}
	chain := GetChain(env, nodeID)
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

func TestSendTransaction(t *testing.T) {
	d := newData(true)
	assert.Equal(t, 50, d.chain.GetBalance(d.address1))
	d.chain.Send(d.wallet, d.address2, 10)
	assert.Equal(t, 40, d.chain.GetBalance(d.address1))
	assert.Equal(t, 10, d.chain.GetBalance(d.address2))
	d.chain.db.Close()
}

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

func TestBlockHeight(t *testing.T) {
	d := newData(true)
	tx, _ := d.chain.NewTransaction(d.address1, d.address2, 1)
	d.chain.SignTransaction(&d.wallet.PrivateKey, tx)
	b1, _ := d.chain.AddBlock([]*Transaction{tx})
	b2, _ := d.chain.AddBlock([]*Transaction{tx})
	b3, _ := d.chain.AddBlock([]*Transaction{tx})
	assert.Equal(t, 1, b1.Height)
	assert.Equal(t, 2, b2.Height)
	assert.Equal(t, 3, b3.Height)
	assert.Equal(t, 3, d.chain.bestHeight)
}

func TestFailGetTransacationUnknown(t *testing.T) {
	d := newData(true)
	_, err := d.chain.GetTransaction([]byte{0x00})
	assert.Contains(t, err.Error(), "transaction not found")
	d.chain.db.Close()
}
