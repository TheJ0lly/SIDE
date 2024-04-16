package metadata

import "crypto/sha256"

type MetaData struct {
	mSource      string
	mDestination string
	mAssetName   string
	mHash        [32]byte
}

type MetadataIE struct {
	Source      string   `json:"source"`
	Destination string   `json:"destination"`
	AssetName   string   `json:"asset_name"`
	Hash        [32]byte `json:"hash"`
}

func CreateNewMetaData(source string, destination string, assetName string) *MetaData {
	md := &MetaData{
		mSource:      source,
		mDestination: destination,
		mAssetName:   assetName,
	}

	md.mHash = sha256.Sum256([]byte(md.GetMetaDataString()))

	return md
}

func (md *MetaData) GetSourceName() string {
	return md.mSource
}

func (md *MetaData) GetDestinationName() string {
	return md.mDestination
}

func (md *MetaData) GetAssetName() string {
	return md.mAssetName
}

func (md *MetaData) GetMetaDataString() string {
	return md.mSource + " " + md.mDestination + " " + md.mAssetName
}

func (md *MetaData) GetMetadataHash() [32]byte {
	return md.mHash
}
