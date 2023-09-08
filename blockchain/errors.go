package blockchain

import (
	"github.com/TheJ0lly/GoChain/prettyfmt"
)

type DataTooBig struct {
	Data string
}

func (dtb *DataTooBig) Error() string {
	return prettyfmt.Sprintf("\"%s\" - is too big! Maximum length allowed: %d!", dtb.Data, data_length)
}

type BlockCapacityReached struct{}

func (bcr *BlockCapacityReached) Error() string {
	return "This block is full! Need new block!"
}
