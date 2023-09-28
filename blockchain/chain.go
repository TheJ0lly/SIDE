package blockchain

import (
	"os"

	"github.com/TheJ0lly/GoChain/asset"
	"github.com/TheJ0lly/GoChain/generalerrors"
	"github.com/TheJ0lly/GoChain/hashtree"
	"github.com/TheJ0lly/GoChain/prettyfmt"
)

type BlockChain struct {
	blocks      []*Block
	databaseDir string
	lastBlock   *Block
}

// CreateNewBlockchain This function will initialize a new blockchain, along with its genesis block, so that the blockchain is ready to use.
func CreateNewBlockchain(dbLoc string) (*BlockChain, error) {
	files, err := os.ReadDir(dbLoc)

	if err != nil {
		return nil, &generalerrors.ReadDirFailed{Dir: dbLoc}
	}

	if len(files) > 0 {
		prettyfmt.WarningF("Folder %s contains files! Do you want to delete them?[y\\n]\n", dbLoc)
		var choice string
		prettyfmt.Scanln(&choice)

		if choice != "y" {
			return nil, &generalerrors.BlockchainDBHasItems{Dir: dbLoc}
		}

		err = clearFolder(dbLoc, files)

		if err != nil {
			return nil, err
		}
	}

	bc := &BlockChain{databaseDir: dbLoc}
	bc.blocks = append(bc.blocks, createGenesisBlock())
	bc.lastBlock = bc.blocks[0]

	return bc, nil
}

// AddData - This function will add some data onto the blockchain.
func (bc *BlockChain) AddData(from string, asset *asset.Asset, destination string) {
	var metaData = from + "_" + asset.GetName() + "_" + destination

	b := bc.getProperBlock()

	b.addDataToBlock(metaData)

	b.hTree.Clear()
	hl := b.getMetaDataHashes()
	hashtree.GenerateTree(hl, b.hTree)
	hash := hashtree.GenerateTree(hl, b.hTree)

	b.currHash = hash[:]
}
