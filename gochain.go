package main

import (
	"log"
	"os"

	"github.com/qza/gochain/core"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "gochain"
	app.Usage = "gochain help"
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func genesis() {
	chain := core.InitChain("address")
	chain.AddBlock([]*core.Transaction{})
	chain.AddBlock([]*core.Transaction{})
	chainIt := chain.Iterator()
	for {
		block := chainIt.Next()
		block.Log()
		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}
