package blockchain

import (
	"encoding/json"
	"os"

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

func (bc *BlockChain) Save_State() {
	bcie := BlockChain_IE{Last_Block_Hash: string(bc.last_block.curr_hash), Database_Dir: bc.database_dir}

	bytes_to_write, err := json.Marshal(bcie)

	if err != nil {
		prettyfmt.ErrorF("%s\n", err.Error())
		return
	}

	err = os.WriteFile("./bcs", bytes_to_write, 0666)

	if err != nil {
		prettyfmt.ErrorF("%s\n", err.Error())
		return
	}
}

func (b *Block) save_state(db_dir string) {
	bie := Block_IE{Current_Hash: string(b.curr_hash), Previous_Hash: string(b.prev_hash), Meta_Data: b.meta_data}

	bytes_to_write, err := json.Marshal(bie)

	if err != nil {
		prettyfmt.ErrorF("%s\n", err.Error())
		return
	}

	file_path := prettyfmt.Sprintf("%s\\%s", db_dir, bie.Current_Hash)

	err = os.WriteFile(file_path, bytes_to_write, 0666)

	if err != nil {
		prettyfmt.ErrorF("%s\n", err.Error())
		return
	}
}

func Load_Blockchain() *BlockChain {

	prettyfmt.Print("Looking for save of the blockchain...\n", prettyfmt.BLUE)

	bytes_read, err := os.ReadFile("./bcs")

	if err != nil {
		prettyfmt.Print("There is no save file for the blockchain! Creating the blockchain from Genesis...\n", prettyfmt.RED)
		return nil
	}

	bcie := &BlockChain_IE{}

	err = json.Unmarshal(bytes_read, bcie)

	if err != nil {
		prettyfmt.ErrorF("%s\n", err.Error())
		return nil
	}

	bc_blocks := rebuild_blockchain([]byte(bcie.Last_Block_Hash), bcie.Database_Dir)

	if bc_blocks == nil {
		return nil
	}

	bc := &BlockChain{database_dir: bcie.Database_Dir, last_block: &bc_blocks[0], blocks: bc_blocks}

	return bc
}

func load_block(block_hash []byte, db_dir string) *Block {
	files, err := os.ReadDir(db_dir)

	if err != nil {
		prettyfmt.ErrorF("%s\n", err.Error())
		return nil
	}

	block_exists := false

	for _, f := range files {
		if f.Name() == string(block_hash) {
			block_exists = true
			break
		}
	}

	if !block_exists {
		prettyfmt.Printf("There is no block with the hash \"%s\"! Creating blockchain from Genesis...\n", prettyfmt.RED, block_hash)
		return nil
	}

	bytes_read, err := os.ReadFile(prettyfmt.Sprintf("%s\\%s", db_dir, block_hash))

	if err != nil {
		prettyfmt.ErrorF("%s\n", err.Error())
		return nil
	}

	bie := &Block_IE{}

	err = json.Unmarshal(bytes_read, bie)

	if err != nil {
		prettyfmt.ErrorF("%s\n", err.Error())
		return nil
	}

	new_block := &Block{meta_data: bie.Meta_Data, prev_hash: []byte(bie.Previous_Hash), curr_hash: []byte(bie.Current_Hash)}

	return new_block
}

func rebuild_blockchain(last_block_hash []byte, db_dir string) []Block {
	files, err := os.ReadDir(db_dir)

	if err != nil {
		prettyfmt.ErrorF("%s\n", err.Error())
		return nil
	}

	if len(files) == 0 {
		prettyfmt.Printf("There are no files in the database! Creating blockchain from Genesis...\n", prettyfmt.RED)
		return nil
	}

	blocks := make([]Block, 0)
	for {
		nb := load_block(last_block_hash, db_dir)

		if nb == nil {
			return nil
		}

		blocks = append(blocks, *nb)

		if len(nb.prev_hash) == 0 {
			break
		}

		last_block_hash = nb.prev_hash
	}

	return blocks
}
