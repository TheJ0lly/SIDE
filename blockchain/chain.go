package blockchain

import (
	"os"

	"github.com/TheJ0lly/GoChain/asset"
	"github.com/TheJ0lly/GoChain/generalerrors"
	"github.com/TheJ0lly/GoChain/hashtree"
	"github.com/TheJ0lly/GoChain/prettyfmt"
)

type BlockChain struct {
	blocks       []*Block
	database_dir string
	last_block   *Block
}

// This function will initialize a new blockchain, along with its genesis block, so that the blockchain is ready to use.
func Create_New_BlockChain(db_loc string) (*BlockChain, error) {
	files, err := os.ReadDir(db_loc)

	if err != nil {
		return nil, &generalerrors.ReadDirFailed{Dir: db_loc}
	}

	if len(files) > 0 {
		prettyfmt.WarningF("Folder %s contains files! Do you want to delete them?[y\\n]\n", db_loc)
		var choice string
		prettyfmt.Scanln(&choice)

		if choice != "y" {
			return nil, &generalerrors.BlockchainDBHasItems{Dir: db_loc}
		}

		err = clear_folder(db_loc, files)

		if err != nil {
			return nil, err
		}
	}

	bc := &BlockChain{database_dir: db_loc}
	bc.blocks = append(bc.blocks, create_genesis_block())
	bc.last_block = bc.blocks[0]

	return bc, nil
}

// This function will add some data onto the blockchain.
func (bc *BlockChain) Add_Data(from string, asset *asset.Asset, destination string) {
	var meta_data string = from + "_" + asset.Get_Name() + "_" + destination

	b := bc.get_proper_block()

	b.add_data_to_block(meta_data)

	b.htree.Clear()
	hl := b.get_meta_data_hashes()
	hashtree.Generate_Tree(hl, b.htree)
	hash := hashtree.Generate_Tree(hl, b.htree)

	b.curr_hash = hash[:]
}
