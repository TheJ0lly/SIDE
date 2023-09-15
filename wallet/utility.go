package wallet

import (
	"os"

	"github.com/TheJ0lly/GoChain/asset"
	"github.com/TheJ0lly/GoChain/prettyfmt"
)

// This function
func (w *Wallet) check_asset_exists(asset_name string) bool {
	files, err := os.ReadDir(w.database_dir)

	if err != nil {
		prettyfmt.ErrorF("%s\n", err.Error())
		return true
	}

	for _, file := range files {
		if asset_name == file.Name() {
			return true
		}
	}

	return false
}

// This function will be used when making transactions
func (w *Wallet) get_asset(asset_name string) *asset.Asset {
	if w.check_asset_exists(asset_name) {
		for _, a := range w.assets {
			if a.Get_Name() == asset_name {
				return a
			}
		}
	}

	prettyfmt.ErrorF("No asset with the name \"%s\" can be found in your wallet\n", asset_name)
	return nil
}
