package core

import (
	"encoding/hex"
	"log"

	"github.com/boltdb/bolt"
)

// UtxoStore struct
type UtxoStore struct {
	Chain *Blockchain
}

// FindUtxo finds all utxos for public key hash
func (utxos *UtxoStore) FindUtxo(pubKeyHash []byte) []TxOutput {
	var res []TxOutput
	chain := utxos.Chain
	bucket := []byte(chain.config.GetDbUtxoBucket())
	err := chain.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			outs := DeserializeOutputs(v)
			for _, out := range outs {
				if out.CanOutputBeUnlocked(pubKeyHash) {
					res = append(res, out)
				}
			}
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	return res
}

// FindSpendableOutputs finds and returns unspent outputs to reference in inputs
func (utxos *UtxoStore) FindSpendableOutputs(pubkeyHash []byte, amount int) (int, map[string][]int) {
	total := 0
	unspent := make(map[string][]int)
	db := utxos.Chain.db
	bucket := []byte(utxos.Chain.config.GetDbUtxoBucket())
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			txID := hex.EncodeToString(k)
			outs := DeserializeOutputs(v)
			for outIdx, out := range outs {
				if out.CanOutputBeUnlocked(pubkeyHash) && total < amount {
					total += out.Value
					unspent[txID] = append(unspent[txID], outIdx)
				}
			}
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	return total, unspent
}

// Reindex makes new index of all utxo in the chain
func (utxos *UtxoStore) Reindex() error {
	chain := utxos.Chain
	bucket := []byte(chain.config.GetDbUtxoBucket())
	err := chain.db.Update(func(tx *bolt.Tx) error {
		tx.DeleteBucket(bucket)
		tx.CreateBucket(bucket)
		return nil
	})
	if err != nil {
		panic(err)
	}
	outs := chain.FindUtxo()
	err = chain.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		if b == nil {
			panic("bucket is nil")
		}
		for txid, out := range outs {
			key, err2 := hex.DecodeString(txid)
			if err2 != nil {
				panic(err2)
			}
			err2 = b.Put(key, SerializeOutputs(out))
			if err2 != nil {
				panic(err2)
			}
		}
		return nil
	})
	return err
}

// Update updates utco index with new block
func (utxos *UtxoStore) Update(block *Block) error {
	db := utxos.Chain.db
	name := utxos.Chain.config.GetDbUtxoBucket()
	err := db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(name))
		for _, tx := range block.Transactions {
			// if !tx.IsCoinbase() {
			for _, vin := range tx.Vin {
				var updatedOuts []TxOutput
				outs := DeserializeOutputs(bucket.Get(vin.Txid))
				for oi, o := range outs {
					if oi != vin.Vout {
						updatedOuts = append(updatedOuts, o)
					}
				}
				if len(updatedOuts) == 0 {
					bucket.Delete(vin.Txid)
				} else {
					bucket.Put(vin.Txid, SerializeOutputs(updatedOuts))
				}
				// }
				bucket.Put(tx.ID, SerializeOutputs(tx.Vout))
			}
		}
		return nil
	})
	return err
}
