package wallet

import (
	"fmt"
	"os"

	"github.com/TheJ0lly/GoChain/asset"
)

// This function
func (w *Wallet) checkAssetExists(assetName string) bool {
	files, err := os.ReadDir(w.mDatabaseDir)

	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return true
	}

	for _, file := range files {
		if assetName == file.Name() {
			return true
		}
	}

	return false
}

// This function will be used when making transactions
func (w *Wallet) getAsset(assetName string) *asset.Asset {
	if w.checkAssetExists(assetName) {
		for _, a := range w.mAssets {
			if a.GetName() == assetName {
				return a
			}
		}
	}

	fmt.Printf("No asset with the name \"%s\" can be found in your wallet\n", assetName)
	return nil
}
