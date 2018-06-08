package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/qza/gochain/core"
	"github.com/urfave/cli"
)

func main() {
	err := godotenv.Load()
	env := &core.EnvConfig{}
	if err != nil {
		panic(err)
	}
	app := cli.NewApp()
	app.Name = "gochain"
	app.Usage = "gochain help"
	app.Commands = []cli.Command{
		{
			Name:    "chain",
			Aliases: []string{"c"},
			Usage:   "chain actions",
			Subcommands: []cli.Command{
				{
					Name:  "init",
					Usage: "initializes new blockchain",
					Action: func(c *cli.Context) error {
						os.RemoveAll(env.GetDbFile())
						wstore := core.NewWalletStore(env)
						wstore.Load(env.GetWalletStoreFile())
						wallet := wstore.CreateWallet()
						address := string(wallet.GetAddress())
						core.InitChain(env, address)
						log.Printf("genesis wallet")
						wallet.Log()
						return nil
					},
				},
				{
					Name:  "print",
					Usage: "prints all blocks in the chain",
					Action: func(c *cli.Context) error {
						chain := core.GetChain(env)
						chain.Log()
						log.Println("ok")
						return nil
					},
				},
			},
		},
		{
			Name:    "balance",
			Aliases: []string{"b"},
			Usage:   "get the balance of address",
			Action: func(c *cli.Context) error {
				chain := core.GetChain(env)
				balance := chain.GetBalance(c.Args().First())
				log.Printf("balance: %d", balance)
				return nil
			},
		},
		{
			Name:    "send",
			Aliases: []string{"s"},
			Usage:   "sends the amount to destination address",
			Action: func(c *cli.Context) error {
				chain := core.GetChain(env)
				wstore := core.NewWalletStore(env)
				wstore.Load(env.GetWalletStoreFile())
				from := c.Args().Get(0)
				to := c.Args().Get(1)
				amount, erra := strconv.ParseInt(c.Args().Get(2), 10, 64)
				if erra != nil {
					panic(erra)
				}
				wallet := wstore.GetWallet(from)
				chain.Send(wallet, to, int(amount))
				return nil
			},
		},
		{
			Name:    "wallet",
			Aliases: []string{"w"},
			Usage:   "wallet actions",
			Subcommands: []cli.Command{
				{
					Name:  "new",
					Usage: "initializes new wallet",
					Action: func(c *cli.Context) error {
						wstore := core.NewWalletStore(env)
						wstore.Load(env.GetWalletStoreFile())
						wallet := wstore.CreateWallet()
						log.Printf("wallet created")
						wallet.Log()
						return nil
					},
				},
				{
					Name:  "get",
					Usage: "gets existing wallet",
					Action: func(c *cli.Context) error {
						wstore := core.NewWalletStore(env)
						erro := wstore.Load(env.GetWalletStoreFile())
						if erro != nil {
							panic(erro)
						}
						wallet := wstore.GetWallet(c.Args().Get(0))
						wallet.Log()
						return nil
					},
				},
				{
					Name:  "utxos",
					Usage: "finds utxos",
					Action: func(c *cli.Context) error {
						chain := core.GetChain(env)
						wstore := core.NewWalletStore(env)
						wstore.Load(env.GetWalletStoreFile())
						utxos := &core.UtxoStore{Chain: chain}
						wallet := wstore.GetWallet(c.Args().Get(0))
						address := string(wallet.GetAddress())
						pubKeyHash, _ := core.PubKeyHash(address)
						unspent := utxos.FindUtxo(pubKeyHash)
						fmt.Printf("UTXOs for address: %s\n", address)
						for _, utxo := range unspent {
							utxo.Log()
						}
						return nil
					},
				},
			},
		},
		{
			Name:    "nodes",
			Aliases: []string{"w"},
			Usage:   "nodes actions",
			Subcommands: []cli.Command{
				{
					Name:  "list",
					Usage: "lists running nodes",
					Action: func(c *cli.Context) error {
						log.Printf("running nodes: \n")
						// TODO:
						return nil
					},
				},
				{
					Name:  "start",
					Usage: "starts new node",
					Action: func(c *cli.Context) error {
						port := c.Args().Get(0)
						mode := c.Args().Get(1)
						// TODO:
						log.Printf("%s node started on port %s \n", mode, port)
						return nil
					},
				},
				{
					Name:  "stop",
					Usage: "stops running node",
					Action: func(c *cli.Context) error {
						port := c.Args().Get(0)
						// TODO:
						log.Printf("node %s stopped \n", port)
						return nil
					},
				},
			},
		},
	}
	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
