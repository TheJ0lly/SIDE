package asset

/*
JPEG = FFD8{...}FFD9
PDF = %PDF / 25 50 44 46

	%  P  D  F
*/
const (
	UNKNOWN FileType = iota
	JPEG
	PDF
)

const (
	jpegHeaderValue  = 0xFFD8
	jpegClosureValue = 0xFFD9

	pdfHeaderValue = 0x25504446
)

type FileType uint8

type Asset struct {
	mName     string
	mFileType FileType
	mData     []byte
}

// CreateNewAsset - will create a new asset structure and will return it as a pointer.
func CreateNewAsset(assetName string, ft FileType, data []byte) *Asset {
	return &Asset{mName: assetName, mFileType: ft, mData: data}
}

// DetermineAssetType - This function will determine the type of asset to upload to the wallet.
func DetermineAssetType(data []byte) FileType {

	if determineJPEG(data) {
		return JPEG
	} else if determinePDF(data) {
		return PDF
	}

	return UNKNOWN

}

// GetName - will return the name of the asset as it is saved by the user.
func (a *Asset) GetName() string {
	return a.mName
}
