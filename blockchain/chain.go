package blockchain

import (
	"fmt"
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

	return bc, nil
}

// AddData - will add some data onto the blockchain.
func (bc *BlockChain) AddData(from string, asset *asset.Asset, destination string) {
	var metaData = from + "_" + asset.GetName() + "_" + destination

	b := bc.getProperBlock()

	b.addDataToBlock(metaData)

	b.mHashTree.ClearTree()
	hl := b.getMetaDataHashes()
	hashtree.GenerateTree(hl, b.mHashTree)
	hash := hashtree.GenerateTree(hl, b.mHashTree)

	b.mCurrHash = hash[:]
}
