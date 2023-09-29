package asset

// determineJPEG - will check to see and determine if the asset is of type JPEG.
func determineJPEG(data []byte) bool {
	/*
		JPEG mData form = FFD8{mData}FFD9
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

// determinePDF - will check to see and determine if the asset is of type PDF.
func determinePDF(data []byte) bool {
	/*
		PDF mData form = 25 50 44 46{mData}
		Thus we look for the first 2 bytes to match the exact same hexadecimal value.
	*/
	var pdfHeader = uint32(data[0])<<24 | uint32(data[1])<<16 | uint32(data[2])<<8 | uint32(data[3])

	return pdfHeader == pdfHeaderValue
}
