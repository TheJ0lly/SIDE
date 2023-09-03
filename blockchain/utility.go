package blockchain

import (
	"crypto/sha256"
	"fmt"
)

func generate_hash(data []byte) []byte {
	gen_hash := sha256.Sum256(data)

	converted_hash := fmt.Sprintf("%X", gen_hash)

	return []byte(converted_hash)

}

func get_block_bytes(data []string) []byte {
	var all_bytes []byte

	for _, str := range data {
		for _, char := range str {
			all_bytes = append(all_bytes, byte(char))
		}
	}
	return all_bytes
}

func check_if_genesis(b *Block) bool {
	if len(b.meta_data) == 0 {
		return false
	}
	return b.meta_data[0] == genesis_name
}
