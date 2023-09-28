package blockchain

import (
	"crypto/sha256"

	"github.com/TheJ0lly/GoChain/hashtree"
)

const (
	genesisName       = "GenesisBlockIsHereAndNotAnywhereElseAndDoNotGoLookForTheValueBecauseYouWillNotFindIt"
	blockDataCapacity = 3
)

type Block struct {
	metaData []string
	prevHash []byte
	currHash []byte
	hTree    *hashtree.Tree
}

// This function should be used only once at initialization of a BlockChain instance.
func createGenesisBlock() *Block {
	md := []string{genesisName}

	hash := sha256.Sum256([]byte(md[0]))

	return &Block{metaData: md, prevHash: nil, currHash: hash[:]}
}

// This function will create a new block, and it will return the new block.
func createNewBlock(b *Block) *Block {
	hash := sha256.Sum256(b.currHash)
	return &Block{metaData: nil, prevHash: b.currHash, currHash: hash[:], hTree: &hashtree.Tree{}}
}

// This function will add data to a block if possible, otherwise it will return an error.
func (b *Block) addDataToBlock(data string) {
	b.metaData = append(b.metaData, data)
}

func (b *Block) getMetaDataHashes() [][32]byte {
	var newList [][32]byte

	for _, md := range b.metaData {
		newList = append(newList, sha256.Sum256([]byte(md)))
	}

	return newList
}

func (bc *BlockChain) getProperBlock() *Block {
	if len(bc.lastBlock.metaData) == blockDataCapacity || checkIfGenesis(bc.lastBlock) {
		return createNewBlock(bc.lastBlock)
	} else {
		return bc.lastBlock
	}
}
