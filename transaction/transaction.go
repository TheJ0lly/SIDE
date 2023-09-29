package transaction

import (
	"crypto/rsa"

	"github.com/TheJ0lly/GoChain/asset"
)

type Transaction struct {
	mFromAddr rsa.PublicKey
	mToAddr   rsa.PublicKey
	mData     asset.Asset
}

// InitializeNewTransaction - This function will create a new transaction to be added on the blockchain
func InitializeNewTransaction(data asset.Asset, source rsa.PublicKey, destination rsa.PublicKey) *Transaction {
	return &Transaction{mFromAddr: source, mToAddr: destination, mData: data}
}
