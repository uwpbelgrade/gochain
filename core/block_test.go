package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

import _ "github.com/joho/godotenv/autoload"

func TestSerializeDeserialize(t *testing.T) {

	a1 := "fXfLnpgQQcHPWz9DnE25JQfXK9ZUYyVRX"
	a1Pub := "7705b24142cf04a1731ad53e2e9f19db0daae9809c3b9466cee9265a54695629acf464682a25ffd594cc9ea3b8325e336cdd66e38a89b54e36fe03f035858318"

	a2 := "VkF377DV5khkmXjw35PnWH1frXKc1vVQS"
	a2Pub := "a3c6ee8418a14ed252438dd40f072b3b874742e2839a2a58542432014a6323213cb9208a76f2ce44361456941cee1ca0cf23a5a126e204e7cc7009a50cb69c09"

	txin1 := TxInput{[]byte("tx1"), 0, nil, []byte(a1)}
	txin2 := TxInput{[]byte("tx2"), 0, nil, []byte(a2)}

	txout1 := TxOutput{100, []byte(a1Pub)}
	txout2 := TxOutput{200, []byte(a2Pub)}
	txout3 := TxOutput{300, []byte(a2Pub)}

	t1 := &Transaction{[]byte(nil), []TxInput{txin1, txin2}, []TxOutput{txout1, txout2, txout3}}

	block1 := NewBlock([]*Transaction{}, nil, 1)
	block2 := NewBlock([]*Transaction{t1}, block1.Hash, 2)

	block := Deserialize(block2.Serialize())

	// assert.True(t, reflect.DeepEqual(block, block2), "Block is serialized/deserialized properly")
	assert.Equal(t, block2.PrevBlockHash, block.PrevBlockHash)
	assert.Equal(t, block2.Hash, block.Hash)
	assert.Equal(t, block2.Nonce, block.Nonce)
	assert.Equal(t, block2.Height, block.Height)

	for i, tx := range block.Transactions {
		tx2 := block2.Transactions[i]
		assert.Equal(t, tx.ID, tx2.ID)
		assert.Equal(t, tx.Vin, tx2.Vin)
		assert.Equal(t, tx.Vout, tx2.Vout)
	}

	block.Log()
}

func DemoTransaction() *Transaction {
	txin1 := TxInput{[]byte("tx1"), 0, nil, []byte("script1")}
	txin2 := TxInput{[]byte("tx1"), 0, nil, []byte("script1")}
	txout1 := TxOutput{100, []byte("address1")}
	txout2 := TxOutput{200, []byte("address2")}
	txout3 := TxOutput{300, []byte("address3")}
	return &Transaction{[]byte{}, []TxInput{txin1, txin2}, []TxOutput{txout1, txout2, txout3}}
}
