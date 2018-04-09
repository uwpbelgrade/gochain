package core

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMerkleTree(t *testing.T) {

	data := [][]byte{
		[]byte("data1"),
		[]byte("data2"),
		[]byte("data3"),
	}

	tree := NewMerkleTree(data)
	treeHash := fmt.Sprintf("%x", tree.Root.Data)

	mn1 := NewMerkleNode(nil, nil, data[0])
	mn2 := NewMerkleNode(nil, nil, data[1])
	mn3 := NewMerkleNode(nil, nil, data[2])
	mn4 := NewMerkleNode(nil, nil, data[2])
	mn5 := NewMerkleNode(mn1, mn2, nil)
	mn6 := NewMerkleNode(mn3, mn4, nil)
	mn7 := NewMerkleNode(mn5, mn6, nil)

	rootHash := fmt.Sprintf("%x", mn7.Data)

	assert.Equal(t, rootHash, treeHash, "Hash is correct")
}
