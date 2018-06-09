package core

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/mr-tron/base58/base58"
)

// Blockchain data structure
type Blockchain struct {
	tip        []byte
	db         *bolt.DB
	bestHeight int
	config     Config
	ws         WalletStore
}

// BlockchainIterator iterates over blocks
type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
	config      Config
}

// InitChain makes new blockchain
func InitChain(config Config, address string, nodeID string) *Blockchain {
	var tip []byte
	ws := NewWalletStore(config, nodeID)
	ws.Load(config.GetWalletStoreFile(nodeID))
	db, err := bolt.Open(config.GetDbFile(), 0600, nil)
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(config.GetDbBucket()))
		if b == nil {
			ts := NewCoinbaseTransaction(address, config.GetGenesisData(), config.GetBlockReward())
			gen := NewBlock([]*Transaction{ts}, []byte{}, 0)
			b, err = tx.CreateBucket([]byte(config.GetDbBucket()))
			err = b.Put(gen.Hash, gen.Serialize())
			err = b.Put([]byte("1"), gen.Hash)
			tip = gen.Hash
		} else {
			tip = b.Get([]byte("1"))
		}
		return nil
	})
	chain := &Blockchain{tip, db, 0, config, *ws}
	utxos := UtxoStore{chain}
	utxos.Reindex()
	return chain
}

// GetBestHeight gets the max block height
func GetBestHeight(db *bolt.DB, config Config) int {
	var lastBlock Block
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(config.GetDbBucket()))
		lastHash := bucket.Get([]byte("1"))
		blockBytes := bucket.Get(lastHash)
		lastBlock = *Deserialize(blockBytes)
		return nil
	})
	if err != nil {
		panic(err)
	}
	return lastBlock.Height
}

// GetChain makes new blockchain
func GetChain(config Config, nodeID string) *Blockchain {
	var tip []byte
	db, err := bolt.Open(config.GetDbFile(), 0600, nil)
	if err != nil {
		panic(err)
	}
	ws := NewWalletStore(config, nodeID)
	ws.Load(config.GetWalletStoreFile(nodeID))
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(config.GetDbBucket()))
		tip = b.Get([]byte("1"))
		return nil
	})
	if err != nil {
		panic(err)
	}
	bestHeight := GetBestHeight(db, config)
	return &Blockchain{tip, db, bestHeight, config, *ws}
}

// AddBlock adds given data as new block in chain
func (chain *Blockchain) AddBlock(ts []*Transaction) (*Block, error) {
	var tip []byte
	for _, tx := range ts {
		if chain.VerifyTransaction(tx) != true {
			return nil, fmt.Errorf("invalid transaction found [txid:%x]", tx.ID)
		}
	}
	err := chain.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(chain.config.GetDbBucket()))
		tip = b.Get([]byte("1"))
		return nil
	})
	newBestHeight := GetBestHeight(chain.db, chain.config) + 1
	block := NewBlock(ts, tip, newBestHeight)
	err = chain.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(chain.config.GetDbBucket()))
		err = b.Put(block.Hash, block.Serialize())
		err = b.Put([]byte("1"), block.Hash)
		chain.tip = block.Hash
		return nil
	})
	chain.bestHeight = newBestHeight
	return block, err
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

// FindUtxo gets unspent transactions outputs
func (chain *Blockchain) FindUtxo() map[string][]TxOutput {
	unspent := make(map[string][]TxOutput)
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
				unspent[txID] = append(unspent[txID], out)
			}
			if tx.IsCoinbase() == false {
				for _, in := range tx.Vin {
					inTxID := hex.EncodeToString(in.Txid)
					spent[inTxID] = append(spent[inTxID], in.Vout)
				}
			}
		}
		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
	return unspent
}

// GetBalance gets balance for address
func (chain *Blockchain) GetBalance(address string) int {
	var balance int
	store := &UtxoStore{chain}
	pubKeyHash, _ := base58.Decode(address)
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-AddressChecksumLength]
	for _, utxo := range store.FindUtxo(pubKeyHash) {
		balance += utxo.Value
	}
	return balance
}

// NewTransaction generates new transaction from spendable outputs
func (chain *Blockchain) NewTransaction(from, to string, amount int) (*Transaction, error) {
	var txins []TxInput
	var txous []TxOutput
	store := &UtxoStore{chain}
	pubKeyHash, _ := PubKeyHash(from)
	spendable, outs := store.FindSpendableOutputs(pubKeyHash, amount)
	if spendable < amount {
		return nil, fmt.Errorf("not enough balance")
	}
	wFrom := chain.ws.GetWallet(from)
	if wFrom == nil {
		fmt.Printf("no such wallet")
		return nil, fmt.Errorf("no such wallet %s", wFrom)
	}
	returnable := spendable - amount
	for txi, touts := range outs {
		txid, err := hex.DecodeString(txi)
		if err != nil {
			fmt.Printf("error decoding txi")
			return nil, error(err)
		}
		for _, out := range touts {
			txins = append(txins, TxInput{txid, out, nil, wFrom.PublicKey})
		}
	}
	txous = append(txous, *NewTxOutput(amount, to))
	txous = append(txous, *NewTxOutput(returnable, from))
	tx := &Transaction{nil, txins, txous}
	tx.ID = tx.Hash()
	return tx, nil
}

// GetTransaction gets transaction by id
func (chain *Blockchain) GetTransaction(id []byte) (Transaction, error) {
	it := chain.Iterator()
	for {
		block := it.Next()
		for _, tx := range block.Transactions {
			if bytes.Compare(id, tx.ID) == 0 {
				return *tx, nil
			}
		}
		if block.PrevBlockHash == nil {
			break
		}
	}
	return Transaction{}, errors.New("transaction not found")
}

// SendFromAddress sends transaction
func (chain *Blockchain) SendFromAddress(pk ecdsa.PrivateKey, from string, to string, amount int) error {
	tx, _ := chain.NewTransaction(from, to, amount)
	chain.SignTransaction(&pk, tx)
	block, _ := chain.AddBlock([]*Transaction{tx})
	store := &UtxoStore{chain}
	return store.Update(block)
}

// Send sends transaction
func (chain *Blockchain) Send(wallet *Wallet, to string, amount int) error {
	return chain.SendFromAddress(wallet.PrivateKey, string(wallet.GetAddress()), to, amount)
}

// GetPreviousTransactions gets previous transactions
func (chain *Blockchain) GetPreviousTransactions(tx Transaction) map[string]Transaction {
	ptxs := make(map[string]Transaction)
	for _, vin := range tx.Vin {
		ptx, _ := chain.GetTransaction(vin.Txid)
		txid := hex.EncodeToString(vin.Txid)
		ptxs[txid] = ptx
	}
	return ptxs
}

// SignTransaction signs transaction
func (chain *Blockchain) SignTransaction(pk *ecdsa.PrivateKey, tx *Transaction) {
	previousTxs := chain.GetPreviousTransactions(*tx)
	tx.Sign(pk, previousTxs)
}

// VerifyTransaction verifies transaction
func (chain *Blockchain) VerifyTransaction(tx *Transaction) bool {
	previousTxs := chain.GetPreviousTransactions(*tx)
	return tx.Verify(previousTxs)
}
