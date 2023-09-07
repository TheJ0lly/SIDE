package wallet

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"os"

	"github.com/TheJ0lly/GoChain/asset"
	"github.com/TheJ0lly/GoChain/database"
)

const (
	key_bit_size = 256
)

type Wallet struct {
	username    string
	password    string
	public_key  rsa.PublicKey
	private_key rsa.PrivateKey
	database    database.Database
	assets      []asset.Asset
}

// This function will create a wallet.
func Create_And_Init_Wallet(username string, password string, location string) *Wallet {
	db := database.Initialize_Database(location)

	if db == nil {
		return nil
	}

	priv_key, err := rsa.GenerateKey(rand.Reader, key_bit_size)

	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return nil
	}

	w := &Wallet{username: username, password: password, private_key: *priv_key, public_key: priv_key.PublicKey, assets: nil}

	w.database = *db

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
		fmt.Printf("Successfully added \"%s\" as an asset.\nFormat: JPEG\nSize: %d bytes\n", asset_name, len(file_data))

		w.database.Write_New_File_To_DB(asset_name)
		w.database.Update_File_From_DB(asset_name, file_data)
		w.assets = append(w.assets, *asset.Create_New_Asset(asset_name, asset.JPEG, file_data))
		return true
	case asset.PDF:
		fmt.Printf("Successfully added \"%s\" as an asset.\nFormat: PDF\nSize: %d bytes\n", asset_name, len(file_data))
		w.database.Write_New_File_To_DB(asset_name)
		w.database.Update_File_From_DB(asset_name, file_data)
		w.assets = append(w.assets, *asset.Create_New_Asset(asset_name, asset.JPEG, file_data))
		return true
	default:
		fmt.Printf("Failed to add \"%s\" as an asset.\nError: Unknown format!\n", asset_name)
		return false
	}

}

// This function will remove an asset from the wallet.
func (w *Wallet) Remove_Asset(asset_name string) {
	if w.check_asset_exists(asset_name) {
		asset_path := fmt.Sprintf("%s\\%s", w.database.Get_Database_Dir(), asset_name)
		err := os.Remove(asset_path)

		if err != nil {
			fmt.Printf("%s\n", err.Error())
			return
		}

		fmt.Printf("Asset \"%s\" has been successfully removed from your wallet\n", asset_name)

	} else {
		fmt.Printf("There is no asset with the name - \"%s\"\n", asset_name)
	}
}

// This function
func (w *Wallet) check_asset_exists(asset_name string) bool {
	files, err := os.ReadDir(w.database.Get_Database_Dir())

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
	asset_path := fmt.Sprintf("%s\\%s", w.database.Get_Database_Dir(), asset_name)

	if w.check_asset_exists(asset_path) {
		index := asset.Find_Asset(asset_name, w.assets)

		if index == -1 {
			fmt.Printf("Weird, the asset \"%s\" is not here\n", asset_name)
			return nil
		}

		return &w.assets[index]
	}

	fmt.Printf("No asset with the name \"%s\" can be found in your wallet\n", asset_name)
	return nil
}
