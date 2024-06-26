package blockchain

import (
	"fmt"
	"github.com/TheJ0lly/GoChain/metadata"
	"github.com/TheJ0lly/GoChain/osspecifics"
	"log"
	"os"

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
		return nil, err
	}

	return bc, nil
}

func CreateNewBlockchainFromData(dbLoc string, blocks []*Block) *BlockChain {
	return &BlockChain{
		mBlocks:      blocks,
		mDatabaseDir: dbLoc,
		mLastBlock:   blocks[len(blocks)-1],
	}
}

// AddData - will add some data onto the blockchain.
func (bc *BlockChain) AddData(from string, destination string, asset string) error {
	md := metadata.CreateNewMetaData(from, destination, asset)

	b, blockExists := bc.getProperBlock()
	b.addDataToBlock(md)

	b.mHashTree.ClearTree()
	hl := getMetaDataHashes(b.mMetaData)

	ht := hashtree.GenerateTree(hl)

	fileName := fmt.Sprintf("%X", b.mCurrHash)

	lastBlockOldHash := osspecifics.CreatePath(bc.mDatabaseDir, fileName)
	err := os.Remove(lastBlockOldHash)

	if err != nil && blockExists {
		log.Printf("ERROR: failed to remove block - %s\n", lastBlockOldHash)
		return &generalerrors.RemoveFileFailed{File: lastBlockOldHash}
	}

	b.mCurrHash = ht.RootHash[:]
	err = ExportBlock(bc.mDatabaseDir, b)

	if err != nil {
		log.Printf("ERROR: failed to export block!\n")
		return &generalerrors.FailedExport{Object: "Block"}
	}

	bc.mBlocks = append(bc.mBlocks, b)
	bc.mLastBlock = b

	return nil
}

// GetDBLocation - return the location of the directory where blockchain data is stored
func (bc *BlockChain) GetDBLocation() string {
	return bc.mDatabaseDir
}

func (bc *BlockChain) GetBlocks() []*Block { return bc.mBlocks }

func (bc *BlockChain) GetLastMetaData() *metadata.MetaData {
	return bc.mLastBlock.mMetaData[len(bc.mLastBlock.mMetaData)-1]
}
