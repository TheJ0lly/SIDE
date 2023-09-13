package osspecifics

import "runtime"

var PATH_SEP string

func init() {
	if runtime.GOOS == "windows" {
		PATH_SEP = "\\"
	} else if runtime.GOOS == "linux" {
		PATH_SEP = "/"
	}
}
