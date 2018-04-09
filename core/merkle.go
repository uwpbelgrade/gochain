package core

import "crypto/sha256"

// MerkleTree datastructure
type MerkleTree struct {
	Root *MerkleNode
}

// MerkleNode of MerkleTree
type MerkleNode struct {
	Left  *MerkleNode
	Right *MerkleNode
	Data  []byte
}

// NewMerkleTree creates tree from given data
// On level 0 there is data, then non-leaf nodes are constructed level by level
func NewMerkleTree(data [][]byte) *MerkleTree {
	var nodes []MerkleNode
	if len(data)%2 != 0 {
		data = append(data, data[len(data)-1])
	}
	for _, datum := range data {
		node := NewMerkleNode(nil, nil, datum)
		nodes = append(nodes, *node)
	}
	for i := 0; i < len(data)/2; i++ {
		var level []MerkleNode
		for j := 0; j < len(nodes); j += 2 {
			node := NewMerkleNode(&nodes[j], &nodes[j+1], nil)
			level = append(level, *node)
		}
		nodes = level
	}
	return &MerkleTree{&nodes[0]}
}

// NewMerkleNode creates new MerkleNode.
// Every Merkle leaf node is labelled with the hash of a data block
// Every Merkle non-leaf node is labelled with the cryptographic hash of the labels of its child node
func NewMerkleNode(left, right *MerkleNode, data []byte) *MerkleNode {
	node := &MerkleNode{left, right, []byte{}}
	if left == nil && right == nil {
		hash := sha256.Sum256(data)
		node.Data = hash[:]
	} else {
		prevHashes := append(left.Data, right.Data...)
		hash := sha256.Sum256(prevHashes)
		node.Data = hash[:]
	}
	return node
}
