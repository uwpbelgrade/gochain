package core

import (
	"fmt"

	"github.com/boltdb/bolt"
)

// BlockBucket blocks bucket name
const BlockBucket = "blocks"

// BlockchainDbFile blocks db file name
const BlockchainDbFile = "/tmp/gochain"

// BlockReward for finding POW
const BlockReward = 50

// GenesisCoinbaseData for genesis block coinbase transaction
const GenesisCoinbaseData = "genesis coinbase data"

// Blockchain data structure
type Blockchain struct {
	tip []byte
	db  *bolt.DB
}

// BlockchainIterator iterates over blocks
type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

// NewCoinbaseTransaction creates new coinbase transaction
func NewCoinbaseTransaction(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Rewart %s", to)
	}
	txin := TxInput{[]byte{}, -1, data}
	txout := TxOutput{BlockReward, to}
	tx := &Transaction{[]byte{}, []TxInput{txin}, []TxOutput{txout}}
	tx.GenerateID()
	return tx
}

// InitChain makes new blockchain
func InitChain(address string) *Blockchain {
	var tip []byte
	db, err := bolt.Open(BlockchainDbFile, 0600, nil)
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlockBucket))
		if b == nil {
			ts := NewCoinbaseTransaction(address, GenesisCoinbaseData)
			gen := NewBlock([]*Transaction{ts}, []byte{})
			b, err = tx.CreateBucket([]byte(BlockBucket))
			err = b.Put(gen.Hash, gen.Serialize())
			err = b.Put([]byte("1"), gen.Hash)
			tip = gen.Hash
		} else {
			tip = b.Get([]byte("1"))
		}
		return nil
	})
	return &Blockchain{tip, db}
}

// AddBlock adds given data as new block in chain
func (chain *Blockchain) AddBlock(ts []*Transaction) {
	var tip []byte
	err := chain.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlockBucket))
		tip = b.Get([]byte("1"))
		return nil
	})
	block := NewBlock(ts, tip)
	err = chain.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlockBucket))
		err = b.Put(block.Hash, block.Serialize())
		err = b.Put([]byte("1"), block.Hash)
		chain.tip = block.Hash
		return nil
	})
}

// Get gets block by hash
func (chain *Blockchain) Get(hash []byte) *Block {
	var block *Block
	err := chain.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlockBucket))
		encodedBlock := b.Get(hash)
		block = Deserialize(encodedBlock)
		return nil
	})
	if err != nil {
		panic(err)
	}
	return block
}

// Iterator makes new Blockchain iterator
func (chain *Blockchain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{chain.tip, chain.db}
}

// Next gets the next block from iterator
func (it *BlockchainIterator) Next() *Block {
	var block *Block
	err := it.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlockBucket))
		encodedBlock := b.Get(it.currentHash)
		block = Deserialize(encodedBlock)
		return nil
	})
	if err != nil {
		panic(err)
	}
	it.currentHash = block.PrevBlockHash
	return block
}
