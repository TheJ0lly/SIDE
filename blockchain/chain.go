package blockchain

import (
	"errors"
)

type BlockChain struct {
	blocks       []Block
	database_dir string
	last_block   *Block
}

// This function will initialize a new blockchain, along with its genesis block, so that the blockchain is ready to use.
func Initialize_BlockChain(db_loc string) *BlockChain {

	saved_bc := load_blockchain()

	if saved_bc != nil {
		return saved_bc
	}

	bc := &BlockChain{database_dir: db_loc}
	bc.blocks = append(bc.blocks, *create_genesis_block())
	bc.last_block = &bc.blocks[0]
	bc.last_block.save_state(db_loc)

	return bc
}

// This function will add some data onto the blockchain.
func (bc *BlockChain) Add_Data(data string) error {
	lb := bc.get_last_block()

	err := lb.add_data_to_block(data)

	if err != nil {
		if errors.Is(err, &DataTooBig{}) {
			return err
		} else if errors.Is(err, &BlockCapacityReached{}) {
			nb := create_new_block(*lb)
			nb.add_data_to_block(data)
			nb.save_state(bc.database_dir)
			bc.blocks = append(bc.blocks, *nb)
			bc.last_block = nb
			return nil
		}
	}

	lb.save_state(bc.database_dir)
	return nil
}
