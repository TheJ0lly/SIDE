package blockchain

import (
	"encoding/json"
	"os"

	"github.com/TheJ0lly/GoChain/generalerrors"
	"github.com/TheJ0lly/GoChain/prettyfmt"
)

type BlockChain_IE struct {
	Last_Block_Hash string `json:"LB_HASH"`
	Database_Dir    string `json:"DB_DIR"`
}

type Block_IE struct {
	Current_Hash  string   `json:"CURRENT_HASH"`
	Previous_Hash string   `json:"PREVIOUS_HASH"`
	Meta_Data     []string `json:"META_DATA"`
}

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

func (b *Block) save_state(db_dir string) error {
	bie := Block_IE{Current_Hash: string(b.curr_hash), Previous_Hash: string(b.prev_hash), Meta_Data: b.meta_data}

	bytes_to_write, err := json.Marshal(bie)

	if err != nil {
		return &generalerrors.JSONMarshalFailed{Object: "Block"}
	}

	file_path := prettyfmt.Sprintf("%s/%s", db_dir, bie.Current_Hash)

	err = os.WriteFile(file_path, bytes_to_write, 0666)

	if err != nil {
		return &generalerrors.WriteFileFailed{File: file_path}
	}

	return nil
}

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

	bc := &BlockChain{database_dir: bcie.Database_Dir, last_block: &bc_blocks[0], blocks: bc_blocks}

	return bc, nil
}

func load_block(block_hash []byte, db_dir string) (*Block, error) {
	files, err := os.ReadDir(db_dir)

	if err != nil {
		return nil, &generalerrors.ReadDirFailed{Dir: db_dir}
	}

	block_exists := false

	for _, f := range files {
		if f.Name() == string(block_hash) {
			block_exists = true
			break
		}
	}

	if !block_exists {
		return nil, &generalerrors.BlockMissing{Block_Hash: string(block_hash)}
	}

	block_file := prettyfmt.Sprintf("%s/%s", db_dir, block_hash)

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
