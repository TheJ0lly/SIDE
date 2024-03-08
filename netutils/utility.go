package netutils

import (
	"bufio"
	"github.com/libp2p/go-libp2p/core/network"
	"io"
	"log"
	"strconv"
)

const (
	AssetNotFound    = -1
	FailedConversion = -2
)

// GetNumberFromResponse - this function will convert the response from a []byte to int
func GetNumberFromResponse(resp []byte) int {

	rStr := ConvertBytesToString(resp)

	n, err := strconv.Atoi(rStr)

	if err != nil {
		log.Printf("ERROR: failed to convert from %s to int\n", rStr)
		return FailedConversion
	}

	return n
}

func ConvertBytesToString(b []byte) string {
	var rStr string

	for _, b := range b {
		if b == 0 {
			break
		}

		rStr += string(b)
	}

	return rStr
}

// ReceiveAsset - this function will try to read the bytes of the asset being sent over the netutils.
func ReceiveAsset(s network.Stream, asset string, val int) (bool, []byte) {
	log.Printf("INFO: current node has asset %s\n", asset)
	log.Printf("INFO: creating buffer of capacity - %d\n", val)
	buff := make([]byte, val)

	_, err := s.Write([]byte("READY"))

	if err != nil {
		log.Printf("ERROR: failed to send ready - %s\n", err)
		return false, nil
	}

	log.Printf("INFO: attempting to read %d bytes\n", val)
	n, err := io.ReadFull(bufio.NewReader(s), buff)

	if err != nil {
		log.Printf("ERROR: %s\n", err)
		return false, nil
	}

	log.Printf("INFO: read %d bytes from the connection\n", n)
	if n != val {
		log.Printf("ERROR: read a different amount of bytes than expected\n")
		return false, nil
	}
	return true, buff
}
