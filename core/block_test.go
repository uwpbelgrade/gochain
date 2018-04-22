package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

import _ "github.com/joho/godotenv/autoload"

func TestSerializeDeserialize(t *testing.T) {

	txin1 := TxInput{[]byte("tx1"), 0, "script1"}
	txin2 := TxInput{[]byte("tx1"), 0, "script1"}

	txout1 := TxOutput{100, "address1"}
	txout2 := TxOutput{200, "address2"}
	txout3 := TxOutput{300, "address3"}

	t1 := &Transaction{[]byte(nil), []TxInput{txin1, txin2}, []TxOutput{txout1, txout2, txout3}}

	block1 := NewBlock([]*Transaction{}, nil)
	block2 := NewBlock([]*Transaction{t1}, block1.Hash)

	block := Deserialize(block2.Serialize())

	// assert.True(t, reflect.DeepEqual(block, block2), "Block is serialized/deserialized properly")
	assert.Equal(t, block2.PrevBlockHash, block.PrevBlockHash)
	assert.Equal(t, block2.Hash, block.Hash)
	assert.Equal(t, block2.Nonce, block.Nonce)

	for i, tx := range block.Transactions {
		tx2 := block2.Transactions[i]
		assert.Equal(t, tx.ID, tx2.ID)
		assert.Equal(t, tx.Vin, tx2.Vin)
		assert.Equal(t, tx.Vout, tx2.Vout)
	}
}

func DemoTransaction() *Transaction {
	txin1 := TxInput{[]byte("tx1"), 0, "script1"}
	txin2 := TxInput{[]byte("tx1"), 0, "script1"}
	txout1 := TxOutput{100, "address1"}
	txout2 := TxOutput{200, "address2"}
	txout3 := TxOutput{300, "address3"}
	return &Transaction{[]byte{}, []TxInput{txin1, txin2}, []TxOutput{txout1, txout2, txout3}}
}
