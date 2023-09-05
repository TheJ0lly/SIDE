package blockchain

import (
	"errors"
	"fmt"
)

type BlockChain struct {
	blocks     []Block
	last_index uint64
}

func Initialize_BlockChain() *BlockChain {
	bc := &BlockChain{last_index: 0}
	bc.blocks = append(bc.blocks, *create_genesis_block())

	return bc
}

func (bc *BlockChain) Get_Last_Block() *Block {
	return &bc.blocks[bc.last_index]
}

func (bc *BlockChain) Add_Data(data string) error {
	lb := bc.Get_Last_Block()

	err := lb.add_data_to_block(data)

	if err != nil {
		if errors.Is(err, &DataTooBig{}) {
			return err
		} else if errors.Is(err, &BlockCapacityReached{}) {
			nb := create_new_block(*lb)
			nb.add_data_to_block(data)
			bc.blocks = append(bc.blocks, *nb)
			bc.last_index++
		}
	}

	return nil
}

// This function is for testing purposes and will soon become obsolete and it will be removed when the database has been introduced
func (bc *BlockChain) Print_BlockChain() {
	for i := uint64(0); i <= bc.last_index; i++ {
		bc.blocks[i].print_block_info()
	}
	fmt.Println()
}
