package blockchain

import (
	"encoding/json"
	"os"

	"github.com/TheJ0lly/GoChain/generalerrors"
	"github.com/TheJ0lly/GoChain/prettyfmt"
)

// This type stands for BlockChain_IMPORT_EXPORT, which will be used to save/load the blockchain state on this machine.
type BlockChain_IE struct {
	Last_Block_Hash string `json:"LB_HASH"`
	Database_Dir    string `json:"DB_DIR"`
}

// This type stands for Block_IMPORT_EXPORT, which will be used to save/load the block state on this machine.
type Block_IE struct {
	Current_Hash  string   `json:"CURRENT_HASH"`
	Previous_Hash string   `json:"PREVIOUS_HASH"`
	Meta_Data     []string `json:"META_DATA"`
}

// This function will save the current state of the blockchain along with its critical info in a JSON file named "bcs", next to the place of execution.
// If it cannot save, it will return an corresponding error.
func (bc *BlockChain) Save_State() error {
	bcie := BlockChain_IE{Last_Block_Hash: string(bc.last_block.curr_hash), Database_Dir: bc.database_dir}

	bytes_to_write, err := json.Marshal(bcie)

	if err != nil {
		return &generalerrors.JSONMarshalFailed{Object: "BlockChain"}
	}

	err = os.WriteFile("./bcs", bytes_to_write, 0666)

	if err != nil {
		return &generalerrors.WriteFileFailed{File: "./bcs"}
	}

	return nil
}

// This function will save the current state of a block in a JSON file having the name of the block's hash, in the folder chosen to be the database.
// If it cannot save, it will return an corresponding error.
func (b *Block) save_state(db_dir string) error {
	bie := Block_IE{Current_Hash: string(b.curr_hash), Previous_Hash: string(b.prev_hash), Meta_Data: b.meta_data}

	bytes_to_write, err := json.Marshal(bie)

	if err != nil {
		return &generalerrors.JSONMarshalFailed{Object: "Block"}
	}

	file_path := prettyfmt.SPathF(db_dir, bie.Current_Hash)

	err = os.WriteFile(file_path, bytes_to_write, 0666)

	if err != nil {
		return &generalerrors.WriteFileFailed{File: file_path}
	}

	return nil
}

// This function will try to load the previous saved state of the blockchain.
// Otherwise, it will return nil and the corresponding error.
func Load_Blockchain() (*BlockChain, error) {

	prettyfmt.CPrint("Looking for save of the blockchain...\n", prettyfmt.BLUE)

	bytes_read, err := os.ReadFile("./bcs")

	if err != nil {
		return nil, &generalerrors.ReadFileFailed{File: "./bcs"}
	}

	bcie := &BlockChain_IE{}

	err = json.Unmarshal(bytes_read, bcie)

	if err != nil {
		return nil, &generalerrors.JSONUnMarshalFailed{Object: "BlockChain"}
	}

	bc_blocks, err := rebuild_blockchain([]byte(bcie.Last_Block_Hash), bcie.Database_Dir)

	if err != nil {
		return nil, err
	}

	err = remove_excess_files(bcie.Database_Dir, bc_blocks)

	if err != nil {
		return nil, err
	}

	bc := &BlockChain{database_dir: bcie.Database_Dir, last_block: &bc_blocks[0], blocks: bc_blocks}

	return bc, nil
}

// This function will recreate the block from memory and add it to the current blockchain instance.
func load_block(block_hash []byte, db_dir string) (*Block, error) {
	block_hash_str := prettyfmt.Sprintf("%s", block_hash)

	block_file := prettyfmt.SPathF(db_dir, block_hash_str)

	_, err := os.Stat(block_file)

	if err != nil {
		return nil, &generalerrors.BlockMissing{Block_Hash: string(block_hash)}
	}

	bytes_read, err := os.ReadFile(block_file)

	if err != nil {
		return nil, &generalerrors.ReadFileFailed{File: block_file}
	}

	bie := &Block_IE{}

	err = json.Unmarshal(bytes_read, bie)

	if err != nil {
		return nil, &generalerrors.JSONUnMarshalFailed{Object: "Block"}
	}

	new_block := &Block{meta_data: bie.Meta_Data, prev_hash: []byte(bie.Previous_Hash), curr_hash: []byte(bie.Current_Hash)}

	return new_block, nil
}

// This function will recreate all the blocks, using only the last hash, and directory of the database.
func rebuild_blockchain(last_block_hash []byte, db_dir string) ([]Block, error) {
	files, err := os.ReadDir(db_dir)

	if err != nil {
		return nil, &generalerrors.ReadDirFailed{Dir: db_dir}
	}

	if len(files) == 0 {
		return nil, &generalerrors.BlockChainDBEmpty{Dir: db_dir}
	}

	blocks := make([]Block, 0)
	for {
		nb, err := load_block(last_block_hash, db_dir)

		if err != nil {
			return nil, err
		}

		blocks = append(blocks, *nb)

		if len(nb.prev_hash) == 0 {
			break
		}

		last_block_hash = nb.prev_hash
	}

	return blocks, nil
}

// This function will remove any files that cannot be reconstructed from the save of the blockchain.
func remove_excess_files(db_dir string, blocks []Block) error {
	files, err := os.ReadDir(db_dir)

	if err != nil {
		return &generalerrors.ReadDirFailed{Dir: db_dir}
	}

	var valid_blocks []string

	for _, b := range blocks {
		valid_blocks = append(valid_blocks, string(b.curr_hash))
	}

	// var len_diff int = len(files) - len(valid_blocks)

	for _, f := range files {
		var to_delete bool = true
		for _, b := range valid_blocks {
			if f.Name() == b {
				to_delete = false
				break
			}
		}

		if to_delete {
			err = os.Remove(prettyfmt.SPathF(db_dir, f.Name()))

			if err != nil {
				return &generalerrors.RemoveFileFailed{File: f.Name()}
			}
		}
	}

	return nil
}
