gochain
-------

Simplified blockchain written in Go with all main functionalities and simplified networking. It helped us get more experience with Go and get better understanding of some internals in how blockchain implementations work.

This example was inspired by original (Jeiwan implementation)[https://github.com/Jeiwan/blockchain_go]. Basic features are same and there are differences internally in how is code structures, environment configuration and command line actions.

## Usage

Use `gochain` command tool to interact with chain.
```
./gochain help
```

First create 2 wallet addresses that will be used by 2 different nodes started on ports `3000` and `3001`.
```
./gochain wallet new 3000
./gochain wallet new 3001
```
This should give you 2 addresses where the block rewards will be sent for the mined blocks.

Now you can start root node:
```
./gochain nodes start 3000 someaddress
```
This will initialize the db and wallet file for node `3000`. You can stop the node.

Now make some transactions using node `3000`.
```
./gochain send 3000 someaddress someaddress2 1
./gochain send 3000 someaddress someaddress2 2
```

Start second node and watch blocks syncing:
```
./gochain nodes start 3001 miner someaddress
```
