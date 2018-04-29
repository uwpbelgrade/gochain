package core

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"

	"github.com/mr-tron/base58/base58"
)

// Version protocol version
const Version = byte(0x01)

// AddressChecksumLength checksum length
const AddressChecksumLength = 4

// Wallet struct
type Wallet struct {
	PrivateKey []byte
	PublicKey  []byte
}

// NewWallet creates new wallet
func NewWallet() *Wallet {
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		panic(err)
	}
	privateK := private.D.Bytes()
	publicK := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)
	wallet := &Wallet{privateK, publicK}
	return wallet
}

// GetAddress gets the wallet address
func (wallet *Wallet) GetAddress() []byte {
	publicK := RipeMd160Sha256(wallet.PublicKey)
	versionedPublicK := append([]byte{Version}, publicK...)
	checksum := ShaChecksum(versionedPublicK, AddressChecksumLength)
	address := base58.Encode(append(versionedPublicK, checksum...))
	return []byte(address)
}

// IsValidAddress validates wallet address
func (wallet *Wallet) IsValidAddress() bool {
	address := wallet.GetAddress()
	publicK, err := base58.Decode(string(address))
	if err != nil {
		panic(err)
	}
	checksum := publicK[len(publicK)-AddressChecksumLength:]
	publicKeyHash := publicK[0 : len(publicK)-AddressChecksumLength]
	requiredChecksum := ShaChecksum(publicKeyHash, AddressChecksumLength)
	return bytes.Compare(checksum, requiredChecksum) == 0
}

// Log prints block info
func (wallet *Wallet) Log() {
	template := `
	WALLET >>>>
	Address: %x
	Public key: %x
	Private key: %x
	`
	fmt.Printf(template, wallet.GetAddress(), wallet.PublicKey, wallet.PrivateKey)
}
