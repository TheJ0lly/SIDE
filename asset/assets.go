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
	name string
	ft   FileType
	data []byte
}

// CreateNewAsset - This function will create a new asset structure and will return it as a pointer.
func CreateNewAsset(assetName string, ft FileType, data []byte) *Asset {
	return &Asset{name: assetName, ft: ft, data: data}
}

// DetermineAssetType - This function will determine the type of asset to upload to the wallet.
func DetermineAssetType(data []byte) FileType {

	if determineJpeg(data) {
		return JPEG
	} else if determinePdf(data) {
		return PDF
	}

	return UNKNOWN

}

func (a *Asset) GetName() string {
	return a.name
}

// This function will check to see and determine if the asset is of type JPEG.
func determineJpeg(data []byte) bool {
	/*
		JPEG data form = FFD8{data}FFD9
		Thus we look for the first and last bytes to match the exact same hexadecimal value.
	*/

	var jpegHeader = uint16(data[0])<<8 | uint16(data[1])

	if jpegHeader != jpegHeaderValue {
		return false
	}

	var dataLen = len(data)

	var jpegClosure = uint16(data[dataLen-2])<<8 | uint16(data[dataLen-1])

	return jpegClosure == jpegClosureValue
}

// This function will check to see and determine if the asset is of type PDF.
func determinePdf(data []byte) bool {
	/*
		PDF data form = 25 50 44 46{data}
		Thus we look for the first 2 bytes to match the exact same hexadecimal value.
	*/
	var pdfHeader = uint32(data[0])<<24 | uint32(data[1])<<16 | uint32(data[2])<<8 | uint32(data[3])

	return pdfHeader == pdfHeaderValue
}
