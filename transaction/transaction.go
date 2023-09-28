package transaction

import (
	"crypto/rsa"

	"github.com/TheJ0lly/GoChain/asset"
)

type Transaction struct {
	fromAddr rsa.PublicKey
	toAddr   rsa.PublicKey
	data     asset.Asset
}

// InitializeNewTransaction - This function will create a new transaction to be added on the blockchain
func InitializeNewTransaction(data asset.Asset, source rsa.PublicKey, destination rsa.PublicKey) *Transaction {
	return &Transaction{fromAddr: source, toAddr: destination, data: data}
}
