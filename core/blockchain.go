package core

import "github.com/boltdb/bolt"

// BlockBucket blocks bucket name
const BlockBucket = "blocks"

// BlockchainDbFile blocks db file name
const BlockchainDbFile = "/tmp/gochain"

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

// InitChain makes new blockchain
func InitChain() *Blockchain {
	var tip []byte
	db, err := bolt.Open(BlockchainDbFile, 0600, nil)
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlockBucket))
		if b == nil {
			gen := NewBlock("genesis", []byte{})
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
func (chain *Blockchain) AddBlock(data string) {
	var tip []byte
	err := chain.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlockBucket))
		tip = b.Get([]byte("1"))
		return nil
	})
	block := NewBlock(data, tip)
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
