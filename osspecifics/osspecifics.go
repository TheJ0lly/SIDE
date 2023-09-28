package osspecifics

import "runtime"

var PathSep string

func init() {
	if runtime.GOOS == "windows" {
		PathSep = "\\"
	} else if runtime.GOOS == "linux" {
		PathSep = "/"
	}
}
