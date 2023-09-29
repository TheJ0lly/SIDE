package blockchain

import (
	"fmt"
	"github.com/TheJ0lly/GoChain/osspecifics"
	"io/fs"
	"os"

	"github.com/TheJ0lly/GoChain/asset"
	"github.com/TheJ0lly/GoChain/generalerrors"
	"github.com/TheJ0lly/GoChain/hashtree"
)

type BlockChain struct {
	mBlocks      []*Block
	mDatabaseDir string
	mLastBlock   *Block
}

// CreateNewBlockchain - will initialize a new blockchain instance.
func CreateNewBlockchain(dbLoc string) (*BlockChain, error) {
	files, err := os.ReadDir(dbLoc)

	if err != nil {
		return nil, &generalerrors.ReadDirFailed{Dir: dbLoc}
	}

	if len(files) > 0 {
		fmt.Printf("Folder %s contains files! Do you want to delete them?[y\\n]\n", dbLoc)
		var choice string
		_, err := fmt.Scanln(&choice)
		if err != nil {
			return nil, err
		}

		if choice != "y" {
			return nil, &generalerrors.BlockchainDBHasItems{Dir: dbLoc}
		}

		err = clearFolder(dbLoc, files)

		if err != nil {
			return nil, err
		}
	}

	bc := &BlockChain{mDatabaseDir: dbLoc}
	bc.mBlocks = append(bc.mBlocks, createGenesisBlock())
	bc.mLastBlock = bc.mBlocks[0]
	err = bc.mLastBlock.Export(dbLoc)

	if err != nil {
		generalerrors.HandleError(err, &generalerrors.AllErrorsExit{ExitCode: 1})
	}

	return bc, nil
}

func clearFolder(dbLoc string, files []fs.DirEntry) error {
	for _, f := range files {
		fileName := osspecifics.CreatePath(dbLoc, f.Name())
		err := os.Remove(fileName)

		if err != nil {
			return err
		}
	}

	return nil
}

// AddData - will add some data onto the blockchain.
func (bc *BlockChain) AddData(from string, asset *asset.Asset, destination string) {
	md := CreateNewMetaData(from, destination, asset.GetName())

	b, blockExists := bc.getProperBlock()
	b.addDataToBlock(md)

	b.mHashTree.ClearTree()
	hl := getMetaDataHashes(b.mMetaData)
	hash := hashtree.GenerateTree(hl, b.mHashTree)

	fileName := fmt.Sprintf("%X", b.mCurrHash)

	err := os.Remove(osspecifics.CreatePath(bc.mDatabaseDir, fileName))

	if err != nil {
		if blockExists {
			fmt.Printf("Error: Failed to remove file!\n")
			return
		}
	}

	b.mCurrHash = hash[:]
	err = b.Export(bc.mDatabaseDir)

	if err != nil {
		fmt.Printf("Error: Failed to export block!\n")
	}
}
