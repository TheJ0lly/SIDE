package blockchain

import (
	"fmt"
	"github.com/TheJ0lly/GoChain/metadata"
	"github.com/TheJ0lly/GoChain/osspecifics"
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

	bc := &BlockChain{mDatabaseDir: dbLoc}
	bc.mBlocks = append(bc.mBlocks, createGenesisBlock())
	bc.mLastBlock = bc.mBlocks[0]
	err := ExportBlock(dbLoc, bc.mLastBlock)

	if err != nil {
		generalerrors.HandleError(err, &generalerrors.AllErrorsExit{ExitCode: 1})
	}

	return bc, nil
}

// AddData - will add some data onto the blockchain.
func (bc *BlockChain) AddData(from string, asset *asset.Asset, destination string) {
	md := metadata.CreateNewMetaData(from, destination, asset.GetName())

	b, blockExists := bc.getProperBlock()
	b.addDataToBlock(md)

	b.mHashTree.ClearTree()
	hl := getMetaDataHashes(b.mMetaData)
	hash := hashtree.GenerateTree(hl, b.mHashTree)

	fmt.Printf("New block hash --- %X\n", hash)

	fileName := fmt.Sprintf("%X", b.mCurrHash)

	err := os.Remove(osspecifics.CreatePath(bc.mDatabaseDir, fileName))

	if err != nil {
		if blockExists {
			fmt.Printf("Error: Failed to remove file!\n")
			return
		}
	}

	b.mCurrHash = hash[:]
	err = ExportBlock(bc.mDatabaseDir, b)

	if err != nil {
		fmt.Printf("Error: Failed to export block!\n")
	}
}
