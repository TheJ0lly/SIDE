package blockchain

import (
	"errors"
	"fmt"

	"github.com/TheJ0lly/GoChain/database"
)

type BlockChain struct {
	blocks     []Block
	database   database.Database
	last_index uint64
}

// This function will initialize a new blockchain, along with its genesis block, so that the blockchain is ready to use.
func Initialize_BlockChain(db_loc string) *BlockChain {
	bc := &BlockChain{last_index: 0}
	bc.blocks = append(bc.blocks, *create_genesis_block())
	bc.database = *database.Initialize_Database(db_loc)
	bc.database.Write_New_File_To_DB(string(bc.get_last_block().curr_hash))
	bc.database.Update_File_From_DB(string(bc.get_last_block().curr_hash), get_block_bytes(bc.get_last_block().meta_data))

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
			bc.blocks = append(bc.blocks, *nb)
			bc.last_index++
			bc.database.Write_New_File_To_DB(string(nb.curr_hash))
			bc.database.Update_File_From_DB(string(nb.curr_hash), get_block_bytes(nb.meta_data))
			return nil
		}
	}
	bc.database.Update_File_From_DB(string(lb.curr_hash), get_block_bytes(lb.meta_data))
	return nil
}

// This function is for testing purposes and will soon become obsolete and it will be removed when the database has been introduced
func (bc *BlockChain) Print_BlockChain() {
	for i := uint64(0); i <= bc.last_index; i++ {
		bc.blocks[i].print_block_info()
	}
	fmt.Println()
}

func (bc *BlockChain) Display_Block_Hashes() {
	fmt.Printf("%s --- Genesis Block\n", bc.blocks[0].curr_hash)

	for i := uint64(1); i <= bc.last_index; i++ {
		fmt.Printf("%s\n", bc.blocks[i].curr_hash)
	}

	fmt.Printf("===== END OF BLOCKS =====\n\n")
}

func (bc *BlockChain) Display_Data_Under_Block(block_hash string) {
	bytes_read, err := bc.database.Read_File_From_DB(block_hash)

	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return
	}

	fmt.Printf("%s\n", bytes_read)
}
