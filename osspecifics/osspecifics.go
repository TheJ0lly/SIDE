package osspecifics

import (
	"github.com/TheJ0lly/GoChain/generalerrors"
	"os"
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

func ClearFolder(folder string) error {
	files, err := os.ReadDir(folder)

	if err != nil {
		return &generalerrors.ReadDirFailed{Dir: folder}
	}

	for _, f := range files {
		path := CreatePath(folder, f.Name())
		err = os.Remove(path)

		if err != nil {
			return &generalerrors.RemoveFileFailed{File: path}
		}
	}

	return nil
}
