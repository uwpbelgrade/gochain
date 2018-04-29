package core

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResetStore(t *testing.T) {
	wstore := makeWalletStore()
	err := os.Remove(wstore.config.GetWalletStoreFile())
	if err != nil {
		panic(err)
	}
	err = wstore.Reset()
	if err != nil {
		panic(err)
	}
	err = wstore.Reset()
	if err != nil {
		panic(err)
	}
}

func TestCreateWallet(t *testing.T) {
	wstore := makeWalletStore()
	wallet := wstore.CreateWallet()
	wstore.Load(wstore.config.GetWalletStoreFile())
	wallet2 := wstore.GetWallet(string(wallet.GetAddress()))
	assert.NotNil(t, wallet)
	assert.NotNil(t, wallet2)
	assert.Equal(t, wallet, wallet2)
}

func TestDeleteWallet(t *testing.T) {
	wstore := makeWalletStore()
	wallet := wstore.CreateWallet()
	address := string(wallet.GetAddress())
	wstore.DeleteWallet(address)
	assert.Nil(t, wstore.GetWallet(address))
}

func makeWalletStore() *WalletStore {
	wallets := make(map[string]*Wallet)
	config := &EnvConfig{}
	wstore := &WalletStore{wallets, config}
	return wstore
}
