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
	jpeg_header_value  = 0xFFD8
	jpeg_closure_value = 0xFFD9

	pdf_header_value = 0x25504446
)

type FileType uint8

type Asset struct {
	name string
	ft   FileType
	data []byte
}

// This function will create a new asset structure and will return it as a pointer.
func Create_New_Asset(asset_name string, ft FileType, data []byte) *Asset {
	return &Asset{name: asset_name, ft: ft, data: data}
}

// This function will determine the type of asset to upload to the wallet.
func Determine_Asset_Type(data []byte) FileType {

	if determine_JPEG(data) {
		return JPEG
	} else if determine_PDF(data) {
		return PDF
	}

	return UNKNOWN

}

func (a *Asset) Get_Name() string {
	return a.name
}

// This function will check to see and determine if the asset is of type JPEG.
func determine_JPEG(data []byte) bool {
	/*
		JPEG data form = FFD8{data}FFD9
		Thus we look for the first and last bytes to match the exact same hexadecimal value.
	*/

	var jpeg_header uint16 = uint16(data[0])<<8 | uint16(data[1])

	if jpeg_header != jpeg_header_value {
		return false
	}

	var data_len int = len(data)

	var jpeg_closure uint16 = uint16(data[data_len-2])<<8 | uint16(data[data_len-1])

	return jpeg_closure == jpeg_closure_value
}

// This function will check to see and determine if the asset is of type PDF.
func determine_PDF(data []byte) bool {
	/*
		PDF data form = 25 50 44 46{data}
		Thus we look for the first 2 bytes to match the exact same hexadecimal value.
	*/
	var pdf_header uint32 = uint32(data[0])<<24 | uint32(data[1])<<16 | uint32(data[2])<<8 | uint32(data[3])

	return pdf_header == pdf_header_value
}
