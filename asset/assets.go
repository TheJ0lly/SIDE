package asset

import "log"

type Asset struct {
	mName string
	mData []byte
}

// CreateNewAsset - will create a new asset structure and will return it as a pointer.
func CreateNewAsset(assetName string, data []byte) *Asset {
	return &Asset{mName: assetName, mData: data}
}

// GetName - will return the name of the asset as it is saved by the user.
func (a *Asset) GetName() string {
	return a.mName
}

func (a *Asset) GetAssetCopy() Asset {
	return Asset{
		mName: a.mName,
		mData: a.mData,
	}
}

// GetAssetBytes - will return the data of the asset
func (a *Asset) GetAssetBytes() []byte {
	return a.mData
}

func (a *Asset) GetAssetSize() int {
	return len(a.mData)
}

func (a *Asset) PrintInfo() {

	log.Printf("Asset info\n  Name: %s\n  Data Length: %d\n", a.mName, len(a.mData))
}
