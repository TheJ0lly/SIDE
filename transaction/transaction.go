package transaction

import (
	"crypto/rsa"

	"github.com/TheJ0lly/GoChain/asset"
)

type Transaction struct {
	from_addr rsa.PublicKey
	to_addr   rsa.PublicKey
	data      asset.Asset
}

// This function will create a new transaction to be added on the blockchain
func Initialize_New_Transaction(data asset.Asset, source rsa.PublicKey, destination rsa.PublicKey) *Transaction {
	return &Transaction{from_addr: source, to_addr: destination, data: data}
}
