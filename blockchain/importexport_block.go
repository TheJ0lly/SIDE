package blockchain

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/TheJ0lly/GoChain/generalerrors"
	"github.com/TheJ0lly/GoChain/hashtree"
	"github.com/TheJ0lly/GoChain/osspecifics"
	"os"
	"strconv"
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

func getByteFromHex(first byte, second byte) byte {
	var final byte

	if first >= '0' && first <= '9' {
		num, _ := strconv.Atoi(string(first))

		final = byte(num) << 4
	} else {
		final = (first - 55) << 4
	}

	if second >= '0' && second <= '9' {
		num, _ := strconv.Atoi(string(second))

		final = final | byte(num)
	} else {
		final = final | (second - 55)
	}

	return final

}

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

func GetMetadataSlice(mie []metadataIE) []*MetaData {
	var mdSlice []*MetaData

	for _, md := range mie {
		mdSlice = append(mdSlice, &MetaData{
			mSource:      md.Source,
			mDestination: md.Destination,
			mAssetName:   md.AssetName,
		})
	}

	return mdSlice
}

func (b *Block) Import(location string) any {
	//UnMarshalling the blockIE
	allBytes, err := os.ReadFile(location)

	if err != nil {
		return &generalerrors.ReadFileFailed{File: location}
	}

	var bie blockIE

	err = json.Unmarshal(allBytes, &bie)

	if err != nil {
		return &generalerrors.JSONUnMarshalFailed{Object: "Block"}
	}

	//Recreating the current hash
	blockHash := osspecifics.GetFileName(location)

	var currentHash [32]byte
	x := 0

	for i := 0; i < len(blockHash); i += 2 {
		currentHash[x] = getByteFromHex(blockHash[i], blockHash[i+1])
		x++
	}

	//Recreating the previous hash
	var previousHash [32]byte
	x = 0

	for i := 0; i < len(bie.PrevHash); i += 2 {
		previousHash[x] = getByteFromHex(bie.PrevHash[i], bie.PrevHash[i+1])
		x++
	}

	//Generating the metadata

	metadata := GetMetadataSlice(bie.MetaData)

	//Generating the hash tree
	ht := &hashtree.Tree{}

	mh := getMetaDataHashes(metadata)

	rootHash := hashtree.GenerateTree(mh, ht)

	if bytes.Compare(currentHash[:], rootHash[:]) != 0 {
		fmt.Printf("Error: Block hash is not equal with the hash tree root")
		return nil
	}

	return &Block{
		mMetaData: metadata,
		mPrevHash: previousHash[:],
		mCurrHash: currentHash[:],
		mHashTree: ht,
	}

}

func (b *Block) Export(folderLocation string) error {
	bie := blockIE{
		MetaData: getMetadataIESlice(b),
		PrevHash: fmt.Sprintf("%X", b.mPrevHash),
	}

	currHashStr := fmt.Sprintf("%X", b.mCurrHash)

	byteToWrite, err := json.Marshal(bie)

	if err != nil {
		return &generalerrors.JSONMarshalFailed{Object: "Block"}
	}

	path := osspecifics.CreatePath(folderLocation, currHashStr)

	err = os.WriteFile(path, byteToWrite, 0666)

	if err != nil {
		return &generalerrors.WriteFileFailed{File: path}
	}

	return nil
}
