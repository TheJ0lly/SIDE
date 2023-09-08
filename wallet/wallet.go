package wallet

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"os"

	"github.com/TheJ0lly/GoChain/asset"
)

const (
	key_bit_size = 256
)

type Wallet struct {
	username     string
	password     string
	public_key   rsa.PublicKey
	private_key  rsa.PrivateKey
	database_dir string
	assets       []*asset.Asset
}

// This function will create a wallet.
func Create_And_Init_Wallet(username string, password string, location string) *Wallet {

	saved_wallet := load_wallet()

	if saved_wallet != nil {
		fmt.Printf("Wallet save found! Restoring...\n")
		return saved_wallet
	}

	priv_key, err := rsa.GenerateKey(rand.Reader, key_bit_size)

	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return nil
	}

	w := &Wallet{username: username, password: password, private_key: *priv_key, public_key: priv_key.PublicKey, assets: nil, database_dir: location}

	return w
}

// This function will add an asset to the wallet.
func (w *Wallet) Add_Asset(asset_name string, file_location string) bool {

	if w.check_asset_exists(asset_name) {
		fmt.Printf("There is already an asset with this name - \"%s\"\n", asset_name)
		return false
	}

	temp, err := os.Stat(file_location)

	if os.IsNotExist(err) {
		fmt.Printf("Asset location does not exist! - \"%s\"\n", file_location)
		return false
	}

	if temp.IsDir() {
		fmt.Printf("Asset is a folder! - \"%s\"\n", file_location)
		return false
	}

	file_data, err := os.ReadFile(file_location)

	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return false
	}

	switch asset.Determine_Asset_Type(file_data) {
	case asset.JPEG:
		asset_to_add := asset.Create_New_Asset(asset_name, asset.JPEG, file_data)
		w.assets = append(w.assets, asset_to_add)
		os.WriteFile(fmt.Sprintf("%s\\%s", w.database_dir, asset_name), file_data, 0444)

		fmt.Printf("Successfully added \"%s\" as an asset.\nFormat: JPEG\nSize: %d bytes\n", asset_name, len(file_data))
		return true
	case asset.PDF:
		asset_to_add := asset.Create_New_Asset(asset_name, asset.PDF, file_data)
		w.assets = append(w.assets, asset_to_add)
		os.WriteFile(fmt.Sprintf("%s\\%s", w.database_dir, asset_name), file_data, 0444)

		fmt.Printf("Successfully added \"%s\" as an asset.\nFormat: PDF\nSize: %d bytes\n", asset_name, len(file_data))
		return true
	default:
		fmt.Printf("Failed to add \"%s\" as an asset.\nError: Unknown format!\n", asset_name)
		return false
	}
}

// This function will remove an asset from the wallet.
func (w *Wallet) Remove_Asset(asset_name string) {
	a := w.get_asset(asset_name)

	if a != nil {
		os.Remove(fmt.Sprintf("%s\\%s", w.database_dir, asset_name))
	}
}

// This function
func (w *Wallet) check_asset_exists(asset_name string) bool {
	files, err := os.ReadDir(w.database_dir)

	if err != nil {
		fmt.Printf("%s\n", err.Error())
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
	asset_path := fmt.Sprintf("%s\\%s", w.database_dir, asset_name)

	if w.check_asset_exists(asset_path) {
		for _, a := range w.assets {
			if a.Get_Name() == asset_name {
				return a
			}
		}
	}

	fmt.Printf("No asset with the name \"%s\" can be found in your wallet\n", asset_name)
	return nil
}
