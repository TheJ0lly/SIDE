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
	mMetaData []string
	mPrevHash []byte
	mCurrHash []byte
	mHashTree *hashtree.Tree
}

// createGenesisBlock - will only be used once at initialization of a BlockChain instance.
func createGenesisBlock() *Block {
	md := []string{genesisName}

	hash := sha256.Sum256([]byte(md[0]))

	return &Block{mMetaData: md, mPrevHash: nil, mCurrHash: hash[:]}
}

// createNewBlock - will create a new block, and it will return the new block.
func createNewBlock(b *Block) *Block {
	hash := sha256.Sum256(b.mCurrHash)
	return &Block{mMetaData: nil, mPrevHash: b.mCurrHash, mCurrHash: hash[:], mHashTree: &hashtree.Tree{}}
}

// addDataToBlock - will add data to a block if possible, otherwise it will return an error.
func (b *Block) addDataToBlock(data string) {
	b.mMetaData = append(b.mMetaData, data)
}

// getMetaDataHashes - will return the hashes of all the metadata in the block.
func (b *Block) getMetaDataHashes() [][32]byte {
	var newList [][32]byte

	for _, md := range b.mMetaData {
		newList = append(newList, sha256.Sum256([]byte(md)))
	}

	return newList
}

// getProperBlock - will return the correct block to add data to.
func (bc *BlockChain) getProperBlock() *Block {
	if len(bc.mLastBlock.mMetaData) == blockDataCapacity || checkIfGenesis(bc.mLastBlock) {
		return createNewBlock(bc.mLastBlock)
	} else {
		return bc.mLastBlock
	}
}

// checkIfGenesis - will check if the block passed is the Genesis block.
func checkIfGenesis(b *Block) bool {
	if len(b.mMetaData) == 0 {
		return false
	}
	return b.mMetaData[0] == genesisName
}
