package core

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
)

// Node struct
type Node struct {
	Host       string
	Port       string
	Address    string
	Env        Config
	Chain      *Blockchain
	MinersAdds string
	Mempool    map[string]Transaction
}

// NewNode creates new node
func NewNode(env Config, port string, minersAddress string) *Node {
	address := net.JoinHostPort(host, port)
	wstore := NewWalletStore(env, port)
	wstore.Load(env.GetWalletStoreFile(port))
	wallet := wstore.CreateWallet()
	coinbaseAddress := string(wallet.GetAddress())
	chain := InitChain(env, coinbaseAddress, port)
	return &Node{host, port, address, env, chain, minersAddress, make(map[string]Transaction)}
}

// Start starts node at specific port server
func (node *Node) Start() {
	listen, err := net.Listen("tcp", node.Address)
	defer listen.Close()
	if err != nil {
		panic(err)
	}
	if node.Address != nodes[0] {
		fmt.Printf("sending version to root node: %s\n", nodes[0])
		node.SendVersionCommand(nodes[0], node.Chain, node.Env)
	}
	for {
		fmt.Printf("server listening on port: %s\n", node.Port)
		conn, err := listen.Accept()
		if err != nil {
			panic(err)
		}
		fmt.Printf("connection established on port: %s\n", node.Port)
		go handleConnection(conn, node, node.Env)
	}
}

func handleConnection(conn net.Conn, node *Node, env Config) {
	payload, err := ioutil.ReadAll(conn)
	if err != nil {
		conn.Close()
		panic(err)
	}
	command := FromBytes(payload[:CommandLength])
	localAddr := conn.LocalAddr().String()
	fmt.Printf("received command: '%s' from: %s \n", command, localAddr)
	switch command {
	case "version":
		node.ReceiveVersionCommand(payload, env, localAddr)
	case "getblocks":
		node.ReceiveGetBlocksCommand(payload, env, localAddr)

	default:
		panic("unknown command")
	}
	conn.Close()
}

// SendVersionCommand handles send version command
func (node *Node) SendVersionCommand(address string, bc *Blockchain, env Config) {
	bestHeight := GetBestHeight(bc.db, env)
	versionCommand := VersionCommand{ProtocolVersion, node.Address, bestHeight}
	payload := EncodeData(versionCommand)
	fmt.Printf(" version command: %x \t %x\n", versionCommand, payload)
	request := append(ToBytes("version"), payload...)
	node.SendData(address, request)
}

// SendGetBlocksCommand handles sending get block command
func (node *Node) SendGetBlocksCommand(address string) {
	payload := EncodeData(GetBlocksCommand{node.Address})
	request := append(ToBytes("getblocks"), payload...)
	node.SendData(address, request)
}

// SendBlock sends block
func (node *Node) SendBlock(address string, b *Block) {
	data := BlocksCommand{address, b.Serialize()}
	payload := EncodeData(data)
	request := append(ToBytes("block"), payload...)
	node.SendData(address, request)
}

// ReceiveVersionCommand handles receiving version command
func (node *Node) ReceiveVersionCommand(request []byte, env Config, remoteAddr string) {
	var buff bytes.Buffer
	var data VersionCommand
	fmt.Printf("processing version command: %s\n", request)
	buff.Write(request[CommandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&data)
	if err != nil {
		log.Panic(err)
	}
	localHeight := GetBestHeight(node.Chain.db, env)
	remoteHeight := data.Height
	fmt.Printf("local vs remote height ::: %d ~ %d\n", localHeight, remoteHeight)
	if localHeight < remoteHeight {
		fmt.Printf("found better height (%d) @ %s\n", remoteHeight, remoteAddr)
		node.SendGetBlocksCommand(data.Origin)
	} else if localHeight > remoteHeight {
		fmt.Printf("found less height (%d) @ %s, sending version\n", remoteHeight, remoteAddr)
		node.SendVersionCommand(data.Origin, node.Chain, env)
	}
	if !KnownNode(data.Origin) {
		fmt.Printf("registering unknown node: %s (%s)\n", data.Origin, remoteAddr)
		nodes = append(nodes, data.Origin)
	}
}

// ReceiveGetBlocksCommand sends inventory on received blocks command
func (node *Node) ReceiveGetBlocksCommand(request []byte, env Config, remoteAddr string) {
	blocks := node.Chain.GetBlockHashes()
	node.SendInventory(remoteAddr, "block", blocks)
}

// SendInventory sends inventory of specific type
func (node *Node) SendInventory(address, kind string, items [][]byte) {
	inventory := InventoryCommand{node.Address, kind, items}
	payload := EncodeData(inventory)
	request := append(ToBytes("inventory"), payload...)
	node.SendData(address, request)
}

// SendData sends data to address
func (node *Node) SendData(address string, data []byte) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Printf("node @ %s is not available\n", address)
		var updatedNodes []string
		for _, node := range nodes {
			if node != address {
				updatedNodes = append(updatedNodes, node)
			}
		}
		nodes = updatedNodes
		return
	}
	defer conn.Close()
	_, err = io.Copy(conn, bytes.NewReader(data))
	if err != nil {
		log.Panic(err)
	}
}
