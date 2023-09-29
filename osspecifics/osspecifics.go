package osspecifics

import (
	"slices"
)

var PathSep string

func CreatePath(format ...string) string {
	var stringToReturn string

	for i := 0; i < len(format); i++ {
		stringToReturn += format[i]

		if i != len(format)-1 {
			stringToReturn += PathSep
		}
	}
	return stringToReturn
}

func GetFileName(filepath string) string {

	//_, filename, ok := strings.Cut(filepath, PathSep)

	var filename []byte

	for i := len(filepath) - 1; i >= 0; i-- {
		if filepath[i] == '\\' {
			break
		} else {
			filename = slices.Insert(filename, 0, filepath[i])
		}
	}

	return string(filename)
}
