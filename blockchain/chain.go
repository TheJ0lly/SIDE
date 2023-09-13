package blockchain

import (
	"errors"
	"os"

	"github.com/TheJ0lly/GoChain/asset"
	"github.com/TheJ0lly/GoChain/generalerrors"
	"github.com/TheJ0lly/GoChain/prettyfmt"
)

type BlockChain struct {
	blocks       []*Block
	database_dir string
	last_block   *Block
}

// This function will initialize a new blockchain, along with its genesis block, so that the blockchain is ready to use.
func Initialize_BlockChain(db_loc string) (*BlockChain, error) {
	files, err := os.ReadDir(db_loc)

	if err != nil {
		return nil, &generalerrors.ReadDirFailed{Dir: db_loc}
	}

	for _, f := range files {
		file_name := prettyfmt.SPathF(db_loc, f.Name())
		err = os.Remove(file_name)

		if err != nil {
			return nil, &generalerrors.RemoveFileFailed{File: file_name}
		}
	}

	bc := &BlockChain{database_dir: db_loc}
	bc.blocks = append(bc.blocks, create_genesis_block())
	bc.last_block = bc.blocks[0]
	bc.last_block.save_state(db_loc)

	return bc, nil
}

// This function will add some data onto the blockchain.
func (bc *BlockChain) Add_Data(data string) error {
	lb := bc.last_block

	err := lb.add_data_to_block(data)

	if err != nil && errors.Is(err, &generalerrors.BlockCapacityReached{Capacity: block_data_capacity}) {
		nb := create_new_block(*lb)
		nb.add_data_to_block(data)
		nb.save_state(bc.database_dir)
		bc.blocks = append(bc.blocks, nb)
		bc.last_block = nb
		return nil
	}

	lb.save_state(bc.database_dir)
	return nil
}

func (bc *BlockChain) Add_Data_Test(from string, asset *asset.Asset, destination string) {
	lb := bc.last_block

	var meta_data string = from + "_" + asset.Get_Name() + "_" + destination

	err := lb.add_data_to_block(meta_data)

	if err != nil {
		nb := create_new_block(*lb)
		nb.add_data_to_block(meta_data)
		nb.save_state(bc.database_dir)
		bc.blocks = append(bc.blocks, nb)
		bc.last_block = nb
		return
	}

	lb.save_state(bc.database_dir)
}

func (bc *BlockChain) View_Blockchain() {
	lb := bc.last_block

	for {
		prettyfmt.Printf("Block Hash: %s\n", lb.curr_hash)
		if !check_if_genesis(lb) {
			prettyfmt.Printf("Meta Data:\n")
			for _, md := range lb.meta_data {
				prettyfmt.Printf("    %s\n", md)
			}
		}

		if string(lb.prev_hash) == "" {
			break
		}

		prettyfmt.Print("\n\n")

		lb = bc.get_block(lb.prev_hash)
	}

}

func (bc *BlockChain) get_block(block_hash []byte) *Block {
	for _, b := range bc.blocks {
		if string(b.curr_hash) == string(block_hash) {
			return b
		}
	}
	return nil
}
