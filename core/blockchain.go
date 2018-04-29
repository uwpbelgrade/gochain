package core

import (
	"encoding/hex"

	"github.com/boltdb/bolt"
)

// Blockchain data structure
type Blockchain struct {
	tip    []byte
	db     *bolt.DB
	config Config
}

// BlockchainIterator iterates over blocks
type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
	config      Config
}

// InitChain makes new blockchain
func InitChain(config Config, address string) *Blockchain {
	var tip []byte
	db, err := bolt.Open(config.GetDbFile(), 0600, nil)
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(config.GetDbBucket()))
		if b == nil {
			ts := NewCoinbaseTransaction(address, config.GetGenesisData(), config.GetBlockReward())
			gen := NewBlock([]*Transaction{ts}, []byte{})
			b, err = tx.CreateBucket([]byte(config.GetDbBucket()))
			err = b.Put(gen.Hash, gen.Serialize())
			err = b.Put([]byte("1"), gen.Hash)
			tip = gen.Hash
		} else {
			tip = b.Get([]byte("1"))
		}
		return nil
	})
	return &Blockchain{tip, db, config}
}

// GetChain makes new blockchain
func GetChain(config Config) *Blockchain {
	var tip []byte
	db, err := bolt.Open(config.GetDbFile(), 0600, nil)
	if err != nil {
		panic(err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(config.GetDbBucket()))
		tip = b.Get([]byte("1"))
		return nil
	})
	if err != nil {
		panic(err)
	}
	return &Blockchain{tip, db, config}
}

// AddBlock adds given data as new block in chain
func (chain *Blockchain) AddBlock(ts []*Transaction) {
	var tip []byte
	err := chain.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(chain.config.GetDbBucket()))
		tip = b.Get([]byte("1"))
		return nil
	})
	block := NewBlock(ts, tip)
	err = chain.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(chain.config.GetDbBucket()))
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
		b := tx.Bucket([]byte(chain.config.GetDbBucket()))
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
	return &BlockchainIterator{chain.tip, chain.db, chain.config}
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
		b := tx.Bucket([]byte(it.config.GetDbBucket()))
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

// GetSpendableOutputs gets the spendable TxOutputs and total amount that covers requested amount
func (chain *Blockchain) GetSpendableOutputs(address string, amount int) (int, map[string][]int) {
	total := 0
	unspent := make(map[string][]int)
	utxs := chain.GetUnspentTransactions(address)
	for utxi := 0; utxi < len(utxs) && total < amount; utxi++ {
		utx := utxs[utxi]
		utxid := hex.EncodeToString(utx.ID)
		for utxoi := 0; utxoi < len(utx.Vout) && total < amount; utxoi++ {
			utxo := utx.Vout[utxoi]
			if utxo.CanOutputBeUnlocked(address) {
				unspent[utxid] = append(unspent[utxid], utxoi)
				total += utxo.Value
			}
		}
	}
	return total, unspent
}

// NewTransaction generates new transaction from spendable outputs
func (chain *Blockchain) NewTransaction(from, to string, amount int) *Transaction {
	var txins []TxInput
	var txous []TxOutput
	spendable, outs := chain.GetSpendableOutputs(from, amount)
	if spendable < amount {
		panic("not enough balance for transaction")
	}
	returnable := spendable - amount
	for txi, touts := range outs {
		txid, err := hex.DecodeString(txi)
		if err != nil {
			panic(err)
		}
		for _, out := range touts {
			txins = append(txins, TxInput{txid, out, from})
		}
	}
	txous = append(txous, TxOutput{amount, to})
	txous = append(txous, TxOutput{returnable, from})
	tx := &Transaction{nil, txins, txous}
	tx.GenerateID()
	return tx
}
