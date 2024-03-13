package blockchain

import (
	"crypto/sha256"
	"github.com/TheJ0lly/GoChain/metadata"

	"github.com/TheJ0lly/GoChain/hashtree"
)

const (
	genesisName       = "GenesisBlockIsHereAndNotAnywhereElseAndDoNotGoLookForTheValueBecauseYouWillNotFindIt"
	blockDataCapacity = 3
)

type Block struct {
	mMetaData []*metadata.MetaData
	mPrevHash []byte
	mCurrHash []byte
	mHashTree *hashtree.Tree
}

// createGenesisBlock - will only be used once at initialization of a BlockChain instance.
func createGenesisBlock() *Block {

	var mdSlice []*metadata.MetaData
	mdSlice = append(mdSlice, metadata.CreateNewMetaData(genesisName, genesisName, genesisName))

	var hashSlice [][32]byte
	hashSlice = append(hashSlice, sha256.Sum256([]byte(mdSlice[0].GetMetaDataString())))

	ht := hashtree.GenerateTree(hashSlice)

	return &Block{
		mMetaData: mdSlice,
		mPrevHash: nil,
		mCurrHash: ht.RootHash[:],
		mHashTree: ht,
	}
}

// createNewBlock - will create a new block, and it will return the new block.
func createNewBlock(b *Block) *Block {
	hash := sha256.Sum256(b.mCurrHash)

	return &Block{
		mMetaData: nil,
		mPrevHash: b.mCurrHash,
		mCurrHash: hash[:],
		mHashTree: &hashtree.Tree{},
	}
}

// addDataToBlock - will add data to a block if possible, otherwise it will return an error.
func (b *Block) addDataToBlock(data *metadata.MetaData) {
	b.mMetaData = append(b.mMetaData, data)
}

// getMetaDataHashes - will return the hashes of all the metadata in the block.
func getMetaDataHashes(md []*metadata.MetaData) [][32]byte {
	var newList [][32]byte

	for _, m := range md {
		newList = append(newList, m.GetMetadataHash())
	}

	return newList
}

// getProperBlock - will return the correct block to add data to.
func (bc *BlockChain) getProperBlock() (*Block, bool) {
	if len(bc.mLastBlock.mMetaData) == blockDataCapacity || checkIfGenesis(bc.mLastBlock) {
		nb := createNewBlock(bc.mLastBlock)
		return nb, false
	} else {
		return bc.mLastBlock, true
	}
}

// checkIfGenesis - will check if the block passed is the Genesis block.
func checkIfGenesis(b *Block) bool {
	if len(b.mMetaData) == 0 {
		return false
	}
	return b.mMetaData[0].GetSourceName() == genesisName && b.mMetaData[0].GetDestinationName() == genesisName && b.mMetaData[0].GetAssetName() == genesisName
}

func (b *Block) GetBlockTreeMatrix() *hashtree.Tree {
	return b.mHashTree
}

func (b *Block) GetMetadata() []*metadata.MetaData {
	return b.mMetaData
}
