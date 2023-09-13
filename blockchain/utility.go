package blockchain

import (
	"crypto/sha256"

	"github.com/TheJ0lly/GoChain/prettyfmt"
)

// This function will generate a sha256 hash from the data given.
func generate_hash(data []byte) []byte {
	gen_hash := sha256.Sum256(data)

	converted_hash := prettyfmt.Sprintf("%X", gen_hash)

	return []byte(converted_hash)

}

// This function will get the block bytes from the meta_data. Soon to be redone when adding []Transactions instead of []string.
func get_block_bytes(data []string) []byte {
	var all_bytes []byte

	for _, str := range data {
		for _, char := range str {
			all_bytes = append(all_bytes, byte(char))
		}
		all_bytes = append(all_bytes, '\n')
	}
	return all_bytes
}

// This function will check if the block passed is the Genesis block.
func check_if_genesis(b *Block) bool {
	if len(b.meta_data) == 0 {
		return false
	}
	return b.meta_data[0] == genesis_name
}
