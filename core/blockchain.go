package core

import (
	"encoding/hex"
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

// GetChain makes new blockchain
func GetChain() *Blockchain {
	var tip []byte
	db, err := bolt.Open(BlockchainDbFile, 0600, nil)
	if err != nil {
		panic(err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlockBucket))
		tip = b.Get([]byte("1"))
		return nil
	})
	if err != nil {
		panic(err)
	}
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

// Log logs current blockchain
func (chain *Blockchain) Log() error {
	chainIt := chain.Iterator()
	for {
		block := chainIt.Next()
		block.Log()
		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
	return nil
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

// GetUnspentTransactions gets unspent transactions
func (chain *Blockchain) GetUnspentTransactions(address string) []Transaction {
	var unspent []Transaction
	spent := make(map[string][]int)
	it := chain.Iterator()
	for {
		block := it.Next()
		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)
		Out:
			for outI, out := range tx.Vout {
				if spent[txID] != nil {
					for _, spentOut := range spent[txID] {
						if spentOut == outI {
							continue Out
						}
					}
				}
				if out.CanOutputBeUnlocked(address) {
					unspent = append(unspent, *tx)
				}
			}
			if tx.IsCoinbase() == false {
				for _, in := range tx.Vin {
					if in.CanUnlockOutput(address) {
						inTxID := hex.EncodeToString(in.Txid)
						spent[inTxID] = append(spent[inTxID], in.Vout)
					}
				}
			}
		}
		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
	return unspent
}

// GetUtxo gets unspent transaction outputs
func (chain *Blockchain) GetUtxo(address string) []TxOutput {
	var unspentOutputs []TxOutput
	for _, tx := range chain.GetUnspentTransactions(address) {
		for _, out := range tx.Vout {
			if out.CanOutputBeUnlocked(address) {
				unspentOutputs = append(unspentOutputs, out)
			}
		}
	}
	return unspentOutputs
}

// GetBalance gets balance for address
func (chain *Blockchain) GetBalance(address string) int {
	var balance int
	for _, utxo := range chain.GetUtxo(address) {
		balance += utxo.Value
	}
	return balance
}
