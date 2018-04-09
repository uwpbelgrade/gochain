package core

// Blockchain data structure
type Blockchain struct {
	Blocks []*Block
}

// AddBlock adds given data as new block in chain
func (chain *Blockchain) AddBlock(data string) {
	last := chain.Blocks[len(chain.Blocks)-1]
	block := NewBlock(data, last.Hash)
	chain.Blocks = append(chain.Blocks, block)
}

// InitChain makes new blockchain
func InitChain() *Blockchain {
	return &Blockchain{[]*Block{NewBlock("genesis", []byte{})}}
}
