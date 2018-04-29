package core

import (
	"crypto/sha256"
	"os"

	"golang.org/x/crypto/ripemd160"
)

// RipeMd160Sha256 ripes md160 sha256
func RipeMd160Sha256(payload []byte) []byte {
	sha := sha256.Sum256(payload)
	riper := ripemd160.New()
	_, err := riper.Write(sha[:])
	if err != nil {
		panic(err)
	}
	ripe := riper.Sum(nil)
	return ripe
}

// ShaChecksum gets the checksum of specific length
func ShaChecksum(payload []byte, length int) []byte {
	sha := sha256.Sum256(payload)
	sha2 := sha256.Sum256(sha[:])
	return sha2[:length]
}

//FileExists checks if file exists on path
func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return !os.IsNotExist(err)
}
