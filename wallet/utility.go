package wallet

import (
	"github.com/TheJ0lly/GoChain/asset"
)

// This function
func (w *Wallet) checkAssetExists(assetName string) bool {
	for _, a := range w.mAssets {
		if a.GetName() == assetName {
			return true
		}
	}

	return false
}

// This function will be used when making transactions
func (w *Wallet) getAsset(assetName string) *asset.Asset {
	for _, a := range w.mAssets {
		if a.GetName() == assetName {
			return a
		}
	}

	return nil
}
