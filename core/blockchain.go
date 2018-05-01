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
	tip    []byte
	db     *bolt.DB
	config Config
	ws     WalletStore
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
	ws := NewWalletStore(config)
	ws.Load(config.GetWalletStoreFile())
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
	return &Blockchain{tip, db, config, *ws}
}

// GetChain makes new blockchain
func GetChain(config Config) *Blockchain {
	var tip []byte
	db, err := bolt.Open(config.GetDbFile(), 0600, nil)
	if err != nil {
		panic(err)
	}
	ws := NewWalletStore(config)
	ws.Load(config.GetWalletStoreFile())
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(config.GetDbBucket()))
		tip = b.Get([]byte("1"))
		return nil
	})
	if err != nil {
		panic(err)
	}
	return &Blockchain{tip, db, config, *ws}
}

// AddBlock adds given data as new block in chain
func (chain *Blockchain) AddBlock(ts []*Transaction) error {
	var tip []byte
	for _, tx := range ts {
		if chain.VerifyTransaction(tx) != true {
			return fmt.Errorf("invalid transaction found [txid:%x]", tx.ID)
		}
	}
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
	return nil
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
func (chain *Blockchain) GetUnspentTransactions(pubKeyHash []byte) []Transaction {
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
				if out.CanOutputBeUnlocked(pubKeyHash) {
					unspent = append(unspent, *tx)
				}
			}
			if tx.IsCoinbase() == false {
				for _, in := range tx.Vin {
					if in.CanUnlockOutput(pubKeyHash) {
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
func (chain *Blockchain) GetUtxo(pubKeyHash []byte) []TxOutput {
	var unspentOutputs []TxOutput
	for _, tx := range chain.GetUnspentTransactions(pubKeyHash) {
		for _, out := range tx.Vout {
			if out.CanOutputBeUnlocked(pubKeyHash) {
				unspentOutputs = append(unspentOutputs, out)
			}
		}
	}
	return unspentOutputs
}

// GetBalance gets balance for address
func (chain *Blockchain) GetBalance(address string) int {
	var balance int
	pubKeyHash, _ := base58.Decode(address)
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-AddressChecksumLength]
	for _, utxo := range chain.GetUtxo(pubKeyHash) {
		balance += utxo.Value
	}
	return balance
}

// GetSpendableOutputs gets the spendable TxOutputs and total amount that covers requested amount
func (chain *Blockchain) GetSpendableOutputs(address string, amount int) (int, map[string][]int) {
	total := 0
	unspent := make(map[string][]int)
	pubKeyHash, _ := base58.Decode(address)
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-AddressChecksumLength]
	utxs := chain.GetUnspentTransactions(pubKeyHash)
	for utxi := 0; utxi < len(utxs) && total < amount; utxi++ {
		utx := utxs[utxi]
		utxid := hex.EncodeToString(utx.ID)
		for utxoi := 0; utxoi < len(utx.Vout) && total < amount; utxoi++ {
			utxo := utx.Vout[utxoi]
			if utxo.CanOutputBeUnlocked(pubKeyHash) {
				unspent[utxid] = append(unspent[utxid], utxoi)
				total += utxo.Value
			}
		}
	}
	return total, unspent
}

// NewTransaction generates new transaction from spendable outputs
func (chain *Blockchain) NewTransaction(from, to string, amount int) (*Transaction, error) {
	var txins []TxInput
	var txous []TxOutput
	spendable, outs := chain.GetSpendableOutputs(from, amount)
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
