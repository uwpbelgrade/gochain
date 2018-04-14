package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSerializeDeserialize(t *testing.T) {

	block1 := NewBlock("data1", nil)
	block2 := NewBlock("data2", block1.Hash)

	block := Deserialize(block2.Serialize())
	assert.Equal(t, block2, block, "Block is serialized/deserialized properly")
}
