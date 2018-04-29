package main

import (
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
						core.InitChain(env, c.Args().First())
						log.Println("ok")
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
				from := c.Args().Get(0)
				to := c.Args().Get(1)
				amount, erra := strconv.ParseInt(c.Args().Get(2), 10, 64)
				if erra != nil {
					panic(erra)
				}
				chain.NewTransaction(from, to, int(amount))
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
						wallet := wstore.CreateWallet()
						log.Printf("wallet created")
						wallet.Log()
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
