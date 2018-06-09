package core

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
)

// CommandLength in bytes
const CommandLength = 12

var nodes []string

// VersionCommand struct
type VersionCommand struct {
	Version int
	Origin  string
	Height  int
}

// SendVersionCommand handles send version command
func SendVersionCommand(address string, bc *Blockchain) {

}

// SendGetBlocksCommand handles sending get block command
func SendGetBlocksCommand(address string) {

}

// ReceiveVersionCommand handles receiving version command
func ReceiveVersionCommand(request []byte, bc *Blockchain, env Config) {
	var buff bytes.Buffer
	var data VersionCommand
	buff.Write(request[CommandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&data)
	if err != nil {
		log.Panic(err)
	}
	myBestHeight := GetBestHeight(bc.db, env)
	foreignerBestHeight := data.Height
	if myBestHeight < foreignerBestHeight {
		SendGetBlocksCommand(data.Origin)
	} else if myBestHeight > foreignerBestHeight {
		SendVersionCommand(data.Origin, bc)
	}
	if !KnownNode(data.Origin) {
		nodes = append(nodes, data.Origin)
	}
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
