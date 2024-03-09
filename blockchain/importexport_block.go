package blockchain

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/TheJ0lly/GoChain/generalerrors"
	"github.com/TheJ0lly/GoChain/hashtree"
	"github.com/TheJ0lly/GoChain/metadata"
	"github.com/TheJ0lly/GoChain/osspecifics"
	"log"
	"os"
	"strconv"
)

const (
	ASCIIToHexDifference = 55
)

type metadataIE struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
	AssetName   string `json:"asset_name"`
}

type blockIE struct {
	MetaData []metadataIE `json:"meta_data"`
	PrevHash string       `json:"prev_hash"`
}

// getByteFromHex - will return a byte formed out of 2 hex characters
func getByteFromHex(first byte, second byte) byte {
	var final byte

	if first >= '0' && first <= '9' {
		num, _ := strconv.Atoi(string(first))

		final = byte(num) << 4
	} else {
		final = (first - ASCIIToHexDifference) << 4
	}

	if second >= '0' && second <= '9' {
		num, _ := strconv.Atoi(string(second))

		final = final | byte(num)
	} else {
		final = final | (second - ASCIIToHexDifference)
	}

	return final

}

// getMetadataIESlice - will return the metadata of transaction in an exportable format
func getMetadataIESlice(b *Block) []metadataIE {
	var mdSlice []metadataIE
	for _, md := range b.mMetaData {
		mdSlice = append(mdSlice, metadataIE{
			Source:      md.GetSourceName(),
			Destination: md.GetDestinationName(),
			AssetName:   md.GetAssetName(),
		})
	}

	return mdSlice
}

// GetMetadataSlice - will generate the metadata.MetaData slice from the exportable format
func GetMetadataSlice(mie []metadataIE) []*metadata.MetaData {
	var mdSlice []*metadata.MetaData

	for _, md := range mie {
		mdSlice = append(mdSlice, metadata.CreateNewMetaData(md.Source, md.Destination, md.AssetName))
	}

	return mdSlice
}

// ImportBlock - will import a single block
func ImportBlock(location string) (*Block, error) {
	//UnMarshalling the blockIE
	allBytes, err := os.ReadFile(location)

	if err != nil {
		return nil, &generalerrors.ReadFileFailed{File: location}
	}

	var bie blockIE

	err = json.Unmarshal(allBytes, &bie)

	if err != nil {
		return nil, err
	}

	//Recreating the current hash
	blockHash := osspecifics.GetFileName(location)

	var currentHash [32]byte
	x := 0

	for i := 0; i < len(blockHash); i += 2 {
		currentHash[x] = getByteFromHex(blockHash[i], blockHash[i+1])
		x++
	}

	//Recreating the previous hash, if exists
	var PrevHash []byte = nil

	if bie.PrevHash != "" {

		var previousHash [32]byte
		x = 0

		for i := 0; i < len(bie.PrevHash); i += 2 {
			previousHash[x] = getByteFromHex(bie.PrevHash[i], bie.PrevHash[i+1])
			x++
		}

		PrevHash = append(PrevHash, previousHash[:]...)
	}

	//Generating the metadata
	md := GetMetadataSlice(bie.MetaData)

	//Generating the hash tree
	mh := getMetaDataHashes(md)

	ht := hashtree.GenerateTree(mh)

	if bytes.Compare(currentHash[:], ht.RootHash[:]) != 0 {
		return nil, &generalerrors.BlockHashDifferent{
			BlockHash:    fmt.Sprintf("%X", currentHash),
			ComputedHash: fmt.Sprintf("%X", ht.RootHash),
		}
	}

	return &Block{
		mMetaData: md,
		mPrevHash: PrevHash,
		mCurrHash: currentHash[:],
		mHashTree: ht,
	}, nil

}

func ImportBlockFromConn(b []byte) *Block {
	var bie blockIE

	err := json.Unmarshal(b, &bie)

	if err != nil {
		log.Printf("ERROR: failed to marshal block - %s\n", err)
		return nil
	}

	md := GetMetadataSlice(bie.MetaData)

	mh := getMetaDataHashes(md)

	ht := hashtree.GenerateTree(mh)

	var PrevHash []byte

	if bie.PrevHash != "" {

		var previousHash [32]byte
		x := 0

		for i := 0; i < len(bie.PrevHash); i += 2 {
			previousHash[x] = getByteFromHex(bie.PrevHash[i], bie.PrevHash[i+1])
			x++
		}

		PrevHash = append(PrevHash, previousHash[:]...)
	}

	return &Block{
		mMetaData: md,
		mPrevHash: PrevHash,
		mCurrHash: ht.RootHash[:],
		mHashTree: ht,
	}
}

func ExportBlockForConn(b *Block) []byte {
	bie := blockIE{
		MetaData: getMetadataIESlice(b),
		PrevHash: fmt.Sprintf("%X", b.mPrevHash),
	}

	byt, err := json.Marshal(bie)

	if err != nil {
		log.Printf("ERROR: %s\n", err)
		return nil
	}

	return byt
}

// ExportBlock - will export a block
func ExportBlock(folderLocation string, b *Block) error {
	bie := blockIE{
		MetaData: getMetadataIESlice(b),
		PrevHash: fmt.Sprintf("%X", b.mPrevHash),
	}

	currHashStr := fmt.Sprintf("%X", b.mCurrHash)

	bytesToWrite, err := json.MarshalIndent(bie, "", "    ")

	if err != nil {
		return err
	}

	path := osspecifics.CreatePath(folderLocation, currHashStr)

	err = os.WriteFile(path, bytesToWrite, 0666)

	if err != nil {
		return &generalerrors.WriteFileFailed{File: path}
	}

	return nil
}
