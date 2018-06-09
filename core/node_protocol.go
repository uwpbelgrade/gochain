package core

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
)

// ProtocolVersion number
const ProtocolVersion = 1.0

// CommandLength in bytes
const CommandLength = 12

var host = "localhost"

var nodes = []string{"localhost:3000"}

// VersionCommand struct
type VersionCommand struct {
	Version int
	Origin  string
	Height  int
}

// GetBlocksCommand struct
type GetBlocksCommand struct {
	Origin string
}

// ToBytes converts command to bytes
func ToBytes(command string) []byte {
	var bytes [CommandLength]byte
	for i, el := range command {
		bytes[i] = byte(el)
	}
	return bytes[:]
}

// FromBytes converts bytes to command string
func FromBytes(bytes []byte) string {
	var command []byte
	for _, el := range bytes {
		if el != 0x0 {
			command = append(command, el)
		}
	}
	return fmt.Sprintf("%s", command)
}

// ExtractCommand extracts command from payload
func ExtractCommand(payload []byte) []byte {
	return payload[:CommandLength]
}

// EncodeData encodes data using gob
func EncodeData(data interface{}) []byte {
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	err := enc.Encode(data)
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}

// KnownNode checks if node is known
func KnownNode(address string) bool {
	for _, node := range nodes {
		if node == address {
			return true
		}
	}
	return false
}
