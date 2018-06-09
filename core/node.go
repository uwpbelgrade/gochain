package core

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
)

// Node struct
type Node struct {
	Host       string
	Port       string
	Chain      *Blockchain
	MinersAdds string
}

// StartNode starts node at specific port server
func StartNode(env Config, port string, minersAddress string) *Node {
	address := net.JoinHostPort("localhost", port)
	listen, err := net.Listen("tcp", address)
	defer listen.Close()
	if err != nil {
		panic(err)
	}
	os.RemoveAll(env.GetDbFile())
	wstore := NewWalletStore(env, port)
	wstore.Load(env.GetWalletStoreFile(port))
	wallet := wstore.CreateWallet()
	coinbaseAddress := string(wallet.GetAddress())
	chain := InitChain(env, coinbaseAddress, port)
	node := &Node{"localhost", port, chain, minersAddress}
	for {
		fmt.Printf("\n server listening on port: %s\n", port)
		conn, err := listen.Accept()
		if err != nil {
			panic(err)
		}
		go handleConnection(conn, chain, env)
		return node
	}
}

func handleConnection(conn net.Conn, chain *Blockchain, env Config) {
	payload, err := ioutil.ReadAll(conn)
	if err != nil {
		conn.Close()
		panic(err)
	}
	command := FromBytes(payload[:CommandLength])
	switch command {
	case "version":
		ReceiveVersionCommand([]byte(command), chain, env)
	default:
		panic("unknown command")
	}
	conn.Close()
}
