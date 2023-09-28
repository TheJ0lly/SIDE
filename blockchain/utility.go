package blockchain

import (
	"io/fs"
	"os"

	"github.com/TheJ0lly/GoChain/prettyfmt"
)

// This function will check if the block passed is the Genesis block.
func checkIfGenesis(b *Block) bool {
	if len(b.metaData) == 0 {
		return false
	}
	return b.metaData[0] == genesisName
}

func clearFolder(dbLoc string, files []fs.DirEntry) error {
	for _, f := range files {
		fileName := prettyfmt.SPathF(dbLoc, f.Name())
		err := os.Remove(fileName)

		if err != nil {
			return err
		}
	}

	return nil
}
