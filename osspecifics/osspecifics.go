package osspecifics

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
