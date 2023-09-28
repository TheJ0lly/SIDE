package blockchain

import (
	"crypto/sha256"

	"github.com/TheJ0lly/GoChain/hashtree"
)

const (
	genesis_name        = "GenesisBlockIsHereAndNotAnywhereElseAndDoNotGoLookForTheValueBecauseYouWillNotFindIt"
	block_data_capacity = 3
)

type Block struct {
	meta_data []string
	prev_hash []byte
	curr_hash []byte
	htree     *hashtree.Tree
}

// This function should be used only once at initialization of a BlockChain instance.
func create_genesis_block() *Block {
	md := []string{genesis_name}

	hash := sha256.Sum256([]byte(md[0]))

	return &Block{meta_data: md, prev_hash: nil, curr_hash: hash[:]}
}

// This function will create a new block, and it will return the new block.
func create_new_block(b *Block) *Block {
	hash := sha256.Sum256(b.curr_hash)
	return &Block{meta_data: nil, prev_hash: b.curr_hash, curr_hash: hash[:], htree: &hashtree.Tree{}}
}

// This function will add data to a block if possible, otherwise it will return an error.
func (b *Block) add_data_to_block(data string) {
	b.meta_data = append(b.meta_data, data)
}

func (b *Block) get_meta_data_hashes() [][32]byte {
	var newList [][32]byte

	for _, md := range b.meta_data {
		newList = append(newList, sha256.Sum256([]byte(md)))
	}

	return newList
}

func (bc *BlockChain) get_proper_block() *Block {
	if len(bc.last_block.meta_data) == block_data_capacity || check_if_genesis(bc.last_block) {
		return create_new_block(bc.last_block)
	} else {
		return bc.last_block
	}
}
