package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/qza/gochain/core"
	"github.com/urfave/cli"
)

func main() {
	err := godotenv.Load()
	env := &core.EnvConfig{}
	if err != nil {
		log.Fatal("Error loading .env file")
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
	}
	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
