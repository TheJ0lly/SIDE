package blockchain

import "fmt"

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

// This function will
func create_new_block(b Block) *Block {
	return &Block{meta_data: nil, prev_hash: b.curr_hash, curr_hash: generate_hash(b.curr_hash)}
}

func (b *Block) add_data_to_block(data string) error {
	if len(data) > data_length {
		return &DataTooBig{Data: data}
	}

	if len(b.meta_data) == block_data_capacity || check_if_genesis(b) {
		return &BlockCapacityReached{}
	}

	b.meta_data = append(b.meta_data, data)
	b.curr_hash = generate_hash(get_block_bytes(b.meta_data))
	return nil
}

func (b *Block) print_block_info() {
	fmt.Printf("previous hash: %s\ncurrent hash: %s\ndata:", b.prev_hash, b.curr_hash)

	for _, md := range b.meta_data {
		fmt.Printf("\n\t%s", md)
	}

	fmt.Printf("\n\n")
}