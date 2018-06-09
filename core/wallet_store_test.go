package core

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResetStore(t *testing.T) {
	wstore := makeWalletStore("1")
	err := os.Remove(wstore.Config.GetWalletStoreFile(wstore.NodeID))
	if err != nil {
		log.Println("wallet file not found")
	}
	err = wstore.Reset(wstore.NodeID)
	if err != nil {
		panic(err)
	}
	err = wstore.Reset(wstore.NodeID)
	if err != nil {
		panic(err)
	}
}

func TestCreateWallet(t *testing.T) {
	wstore := makeWalletStore("1")
	wallet := wstore.CreateWallet()
	wstore.Load(wstore.Config.GetWalletStoreFile(wstore.NodeID))
	wallet2 := wstore.GetWallet(string(wallet.GetAddress()))
	assert.NotNil(t, wallet)
	assert.NotNil(t, wallet2)
	assert.Equal(t, wallet, wallet2)
}

func TestDeleteWallet(t *testing.T) {
	wstore := makeWalletStore("1")
	wallet := wstore.CreateWallet()
	address := string(wallet.GetAddress())
	wstore.DeleteWallet(address)
	assert.Nil(t, wstore.GetWallet(address))
}

func makeWalletStore(nodeID string) *WalletStore {
	wallets := make(map[string]*Wallet)
	config := &EnvConfig{}
	wstore := &WalletStore{wallets, config, nodeID}
	return wstore
}
