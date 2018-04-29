package core

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
)

// WalletStore holds wallets list
type WalletStore struct {
	Wallets map[string]*Wallet
	config  Config
}

// Reset resets the wallet store
func (ws *WalletStore) Reset() error {
	file := ws.config.GetWalletStoreFile()
	if !FileExists(file) {
		ws.Save()
	} else {
		ws.ClearAll()
	}
	err := ws.Load(file)
	if err != nil {
		fmt.Printf("failed loading wallets file")
		return err
	}
	return nil
}

// GetWallet gets the wallet
func (ws *WalletStore) GetWallet(address string) *Wallet {
	return ws.Wallets[address]
}

// CreateWallet creates and saves wallet to file
func (ws *WalletStore) CreateWallet() *Wallet {
	wallet := NewWallet()
	address := string(wallet.GetAddress())
	ws.Wallets[address] = wallet
	ws.Save()
	return wallet
}

// DeleteWallet deletes wallet
func (ws *WalletStore) DeleteWallet(address string) {
	wallet := ws.Wallets[address]
	if wallet == nil {
		log.Println("wallet not found")
		return
	}
	delete(ws.Wallets, address)
	ws.Save()
}

// Save saves wallets to file
func (ws *WalletStore) Save() {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(ws.Wallets)
	if err != nil {
		panic(err)
	}
	file := ws.config.GetWalletStoreFile()
	err = ioutil.WriteFile(file, buffer.Bytes(), 0644)
	if err != nil {
		panic(err)
	}
}

// Load loads wallets from file
func (ws *WalletStore) Load(filename string) error {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	reader := bytes.NewReader(content)
	decoder := gob.NewDecoder(reader)
	err = decoder.Decode(&ws.Wallets)
	if err != nil {
		return err
	}
	return nil
}

// ClearAll deletes all wallets from file
func (ws *WalletStore) ClearAll() {
	ws.Wallets = make(map[string]*Wallet)
	ws.Save()
}
