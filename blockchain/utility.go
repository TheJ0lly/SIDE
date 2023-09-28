package blockchain

import (
	"io/fs"
	"os"

	"github.com/TheJ0lly/GoChain/prettyfmt"
)

// This function will check if the block passed is the Genesis block.
func check_if_genesis(b *Block) bool {
	if len(b.meta_data) == 0 {
		return false
	}
	return b.meta_data[0] == genesis_name
}

func clear_folder(db_loc string, files []fs.DirEntry) error {
	for _, f := range files {
		file_name := prettyfmt.SPathF(db_loc, f.Name())
		err := os.Remove(file_name)

		if err != nil {
			return err
		}
	}

	return nil
}
