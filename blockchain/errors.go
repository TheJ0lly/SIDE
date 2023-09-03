package blockchain

import "fmt"

type DataTooBig struct {
	Data string
}

func (dtb *DataTooBig) Error() string {
	return fmt.Sprintf("\"%s\" - is too big! Maximum length allowed: %d!", dtb.Data, data_length)
}

type BlockCapacityReached struct{}

func (bcr *BlockCapacityReached) Error() string {
	return "This block is full! Need new block!"
}
