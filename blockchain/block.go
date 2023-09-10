package blockchain

import "github.com/TheJ0lly/GoChain/generalerrors"

const (
	genesis_name        = "GenesisBlockIsHereAndNotAnywhereElseAndDoNotGoLookForTheValueBecauseYouWillNotFindIt"
	block_data_capacity = 3
	data_length         = 10
)

type Block struct {
	meta_data []string
	prev_hash []byte
	curr_hash []byte
}

// This function should be used only once at initialization of a BlockChain instance.
func create_genesis_block() *Block {
	md := []string{genesis_name}

	return &Block{meta_data: md, prev_hash: nil, curr_hash: generate_hash(get_block_bytes(md))}
}

// This function will create a new block, and it will return the new block
func create_new_block(b Block) *Block {
	return &Block{meta_data: nil, prev_hash: b.curr_hash, curr_hash: generate_hash(b.curr_hash)}
}

// This function will add data to a block if possible, otherwise it will return an error.
//
//	return values:
//	-DataTooBig - meaning that the data is too big, thus it will require adjustment.
//	-BlockCapacityReached - meaning that the current block capacity has been reached and a new block is needed for the addition and storage of the data.
func (b *Block) add_data_to_block(data string) error {
	if len(data) > data_length {
		return &generalerrors.DataTooBig{Data: data}
	}

	if len(b.meta_data) == block_data_capacity || check_if_genesis(b) {
		return &generalerrors.BlockCapacityReached{}
	}

	b.meta_data = append(b.meta_data, data)
	return nil
}
