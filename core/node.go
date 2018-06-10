package core

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
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
	Transit    [][]byte
}

// NewNode creates new node
func NewNode(env Config, port string, minersAddress string) *Node {
	address := net.JoinHostPort(host, port)
	wstore := NewWalletStore(env, port)
	wstore.Load(env.GetWalletStoreFile(port))
	wallet := wstore.CreateWallet()
	coinbaseAddress := string(wallet.GetAddress())
	chain := InitChain(env, coinbaseAddress, port)
	return &Node{host, port, address, env, chain, minersAddress, make(map[string]Transaction), [][]byte{}}
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
	fmt.Printf("received command: '%s' \n", command)
	switch command {
	case "version":
		node.ReceiveVersionCommand(payload, env)
	case "getblocks":
		node.ReceiveGetBlocksCommand(payload, env)
	case "inventory":
		node.ReceiveInventoryCommand(payload, env)
	case "getdata":
		node.ReceiveGetDataCommand(payload, env)
	case "block":
		node.ReceiveBlockCommand(payload, env)
	case "transaction":
		node.ReceiveTransactionCommand(payload, env)
	default:
		panic("unknown command")
	}
	conn.Close()
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

// SendGetDataCommand sends getdata command
func (node *Node) SendGetDataCommand(address, kind string, id []byte) {
	payload := EncodeData(GetDataCommand{node.Address, kind, id})
	request := append(ToBytes("getdata"), payload...)
	node.SendData(address, request)
}

// SendBlockCommand sends block
func (node *Node) SendBlockCommand(address string, b *Block) {
	data := BlockCommand{node.Address, b.Serialize()}
	payload := EncodeData(data)
	request := append(ToBytes("block"), payload...)
	node.SendData(address, request)
}

// SendTransaction sends transaction
func (node *Node) SendTransaction(address string, t *Transaction) {
	data := TransactionCommand{node.Address, t.Serialize()}
	payload := EncodeData(data)
	request := append(ToBytes("transaction"), payload...)
	node.SendData(address, request)
}

// SendInventory sends inventory of specific type
func (node *Node) SendInventory(address, kind string, items [][]byte) {
	inventory := InventoryCommand{node.Address, kind, items}
	payload := EncodeData(inventory)
	request := append(ToBytes("inventory"), payload...)
	node.SendData(address, request)
}

// ReceiveVersionCommand handles receiving version command
func (node *Node) ReceiveVersionCommand(request []byte, env Config) {
	var buff bytes.Buffer
	var data VersionCommand
	buff.Write(request[CommandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&data)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("processing version command from %s\n", data.Origin)
	localHeight := GetBestHeight(node.Chain.db, env)
	remoteHeight := data.Height
	fmt.Printf("local vs remote height ::: %d ~ %d\n", localHeight, remoteHeight)
	if localHeight < remoteHeight {
		node.SendGetBlocksCommand(data.Origin)
	} else if localHeight > remoteHeight {
		node.SendVersionCommand(data.Origin, node.Chain, env)
	}
	if !KnownNode(data.Origin) {
		fmt.Printf("registering unknown node: %s \n", data.Origin)
		nodes = append(nodes, data.Origin)
	}
}

// ReceiveGetBlocksCommand sends inventory on received blocks command
func (node *Node) ReceiveGetBlocksCommand(request []byte, env Config) {
	var buff bytes.Buffer
	var command GetBlocksCommand
	buff.Write(request[CommandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&command)
	if err != nil {
		log.Panic(err)
	}
	blocks := node.Chain.GetBlockHashes()
	node.SendInventory(command.Origin, "block", blocks)
}

// ReceiveInventoryCommand handles inventory command
func (node *Node) ReceiveInventoryCommand(request []byte, env Config) {
	var buff bytes.Buffer
	var payload InventoryCommand
	buff.Write(request[CommandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("recevied inventory with %d %s\n", len(payload.Data), payload.Type)
	if payload.Type == "block" {
		node.Transit = payload.Data
		blockHash := payload.Data[0]
		node.SendGetDataCommand(payload.Origin, "block", blockHash)
		newInTransit := [][]byte{}
		for _, b := range node.Transit {
			if bytes.Compare(b, blockHash) != 0 {
				newInTransit = append(newInTransit, b)
			}
		}
		fmt.Printf("new in transit: %d \n", len(newInTransit))
		node.Transit = newInTransit
	}
	if payload.Type == "transaction" {
		txID := payload.Data[0]
		if node.Mempool[hex.EncodeToString(txID)].ID == nil {
			node.SendGetDataCommand(payload.Origin, "transaction", txID)
		}
	}
}

// ReceiveBlockCommand processes block command
func (node *Node) ReceiveBlockCommand(request []byte, env Config) {
	var buff bytes.Buffer
	var payload BlockCommand
	buff.Write(request[CommandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}
	blockData := payload.Block
	block := Deserialize(blockData)
	node.Chain.AddBlock(block)
	fmt.Printf("added new block [height: %x] [hash: %x] \n", block.Height, block.Hash)
	if len(node.Transit) > 0 {
		blockHash := node.Transit[0]
		fmt.Printf("fetching next node in transit from %s [hash: %x] \n", payload.Origin, blockHash)
		node.SendGetDataCommand(payload.Origin, "block", blockHash)
		node.Transit = node.Transit[1:]
		fmt.Printf("new transit size %d \n", len(node.Transit))
	} else {
		UTXOSet := UtxoStore{node.Chain}
		UTXOSet.Reindex()
	}
}

// ReceiveTransactionCommand receives transaction command
func (node *Node) ReceiveTransactionCommand(request []byte, env Config) {
	var buff bytes.Buffer
	var payload TransactionCommand
	buff.Write(request[CommandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}
	txData := payload.Transaction
	tx := DeserializeTransaction(txData)
	node.Mempool[hex.EncodeToString(tx.ID)] = tx
	if node.Address == nodes[0] {
		for _, n := range nodes {
			if n != node.Address && n != payload.Origin {
				node.SendInventory(n, "transaction", [][]byte{tx.ID})
			}
		}
	} else {
		if len(node.Mempool) >= 2 && len(node.MinersAdds) > 0 {
		MineTransactions:
			var txs []*Transaction
			for id := range node.Mempool {
				tx := node.Mempool[id]
				if node.Chain.VerifyTransaction(&tx) {
					txs = append(txs, &tx)
				}
			}

			if len(txs) == 0 {
				fmt.Println("all transactions are invalid! waiting for new ones ...")
				return
			}

			cbTx := NewCoinbaseTransaction(node.MinersAdds, "", node.Env.GetBlockReward())
			txs = append(txs, cbTx)

			newBlock, err := node.Chain.MineBlock(txs)
			if err != nil {
				panic(err)
			}
			UTXOSet := UtxoStore{node.Chain}
			UTXOSet.Reindex()

			fmt.Println("new block is mined!")
			for _, tx := range txs {
				txID := hex.EncodeToString(tx.ID)
				delete(node.Mempool, txID)
			}
			for _, n := range nodes {
				if n != node.Address {
					node.SendInventory(n, "block", [][]byte{newBlock.Hash})
				}
			}
			if len(node.Mempool) > 0 {
				goto MineTransactions
			}
		}
	}
}

// ReceiveGetDataCommand handles getdata command
func (node *Node) ReceiveGetDataCommand(request []byte, env Config) {
	var buff bytes.Buffer
	var payload GetDataCommand
	buff.Write(request[CommandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}
	if payload.Type == "block" {
		block, err := node.Chain.GetBlock([]byte(payload.ID))
		if err != nil {
			return
		}
		node.SendBlockCommand(payload.Orign, &block)
	}
	if payload.Type == "transaction" {
		txID := hex.EncodeToString(payload.ID)
		tx := node.Mempool[txID]
		node.SendTransaction(payload.Orign, &tx)
	}
}
