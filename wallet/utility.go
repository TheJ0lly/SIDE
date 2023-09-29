package wallet

import (
	"fmt"
	"github.com/TheJ0lly/GoChain/osspecifics"
	"io/fs"
	"os"

	"github.com/TheJ0lly/GoChain/asset"
)

// This function
func (w *Wallet) checkAssetExists(assetName string) bool {
	files, err := os.ReadDir(w.databaseDir)

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
		for _, a := range w.assets {
			if a.GetName() == assetName {
				return a
			}
		}
	}

	fmt.Printf("No asset with the name \"%s\" can be found in your wallet\n", assetName)
	return nil
}

func clearFolder(dbLoc string, files []fs.DirEntry) error {
	for _, f := range files {
		fileName := osspecifics.CreatePath(dbLoc, f.Name())
		err := os.Remove(fileName)

		if err != nil {
			return err
		}
	}

	return nil
}
