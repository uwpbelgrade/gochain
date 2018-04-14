package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/qza/gochain/core"
)

var addr = flag.String("http", "127.0.0.1", "server listen address")
var port = flag.String("port", "8080", "server listen port")

func boot() {
	flag.Parse()
	log.Printf("\nserver starting")
}

func run() {
	addr := fmt.Sprintf("%s:%s", *addr, *port)
	fmt.Printf("listening on %v \n", addr)
	// http.ListenAndServe(addr, nil)
}

func genesis() {
	chain := core.InitChain()
	chain.AddBlock("block1")
	chain.AddBlock("block2")
	chainIt := chain.Iterator()
	for {
		block := chainIt.Next()
		fmt.Printf("\n")
		fmt.Printf("Previous hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %X\n", block.Hash)
		fmt.Printf("Timestamp: %d [%s]\n", block.Timestamp, time.Unix(block.Timestamp, 0))
		fmt.Printf("Nonce: %s", block.Hash)
		fmt.Printf("\n")
	}
}

func main() {
	boot()
	run()
	genesis()
}
