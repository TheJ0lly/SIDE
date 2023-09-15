package wallet

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"os"

	"github.com/TheJ0lly/GoChain/asset"
	"github.com/TheJ0lly/GoChain/generalerrors"
	"github.com/TheJ0lly/GoChain/prettyfmt"
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
func Initialize_Wallet(username string, password string, db_loc string) (*Wallet, error) {
	files, err := os.ReadDir(db_loc)

	if err != nil {
		prettyfmt.ErrorF("%s\n", err.Error())
		return nil, &generalerrors.ReadDirFailed{Dir: db_loc}
	}

	if len(files) > 0 {
		prettyfmt.WarningF("Folder %s contains files! Do you want to delete them?[y\\n]\n", db_loc)
		var choice string
		prettyfmt.Scanln(&choice)

		if choice != "y" {
			return nil, &generalerrors.WalletDBHasItems{Dir: db_loc}
		}
	}

	for _, f := range files {
		os.Remove(prettyfmt.SPathF(db_loc, f.Name()))
	}

	priv_key, err := rsa.GenerateKey(rand.Reader, key_bit_size)

	if err != nil {
		prettyfmt.ErrorF("%s\n", err.Error())
		return nil, err
	}

	pass_bytes := sha256.Sum256([]byte(password))

	pass_bytes_str := prettyfmt.Sprintf("%X", pass_bytes)

	w := &Wallet{username: username, password: pass_bytes_str, private_key: *priv_key, public_key: priv_key.PublicKey, assets: nil, database_dir: db_loc}

	return w, nil
}

// This function will add an asset to the wallet.
func (w *Wallet) Add_Asset(asset_name string, file_location string) bool {

	if w.check_asset_exists(asset_name) {
		prettyfmt.ErrorF("There is already an asset with this name - \"%s\"\n", asset_name)
		return false
	}

	temp, err := os.Stat(file_location)

	if os.IsNotExist(err) {
		prettyfmt.ErrorF("Asset location does not exist! - \"%s\"\n", file_location)
		return false
	}

	if temp.IsDir() {
		prettyfmt.ErrorF("Asset is a folder! - \"%s\"\n", file_location)
		return false
	}

	file_data, err := os.ReadFile(file_location)

	if err != nil {
		prettyfmt.ErrorF("%s\n", err.Error())
		return false
	}

	switch asset.Determine_Asset_Type(file_data) {
	case asset.JPEG:
		asset_to_add := asset.Create_New_Asset(asset_name, asset.JPEG, file_data)
		w.assets = append(w.assets, asset_to_add)
		err = os.WriteFile(prettyfmt.SPathF(w.database_dir, asset_name), file_data, 0444)

		if err != nil {
			prettyfmt.ErrorF("Failed to add \"%s\" as an asset.\nError: ", asset_name)
			generalerrors.HandleError(err)
			return false
		}

		prettyfmt.CPrintf("Successfully added \"%s\" as an asset.\nFormat: JPEG\nSize: %d bytes\n", prettyfmt.GREEN, asset_name, len(file_data))
		return true
	case asset.PDF:
		asset_to_add := asset.Create_New_Asset(asset_name, asset.PDF, file_data)
		w.assets = append(w.assets, asset_to_add)
		err = os.WriteFile(prettyfmt.SPathF(w.database_dir, asset_name), file_data, 0444)

		if err != nil {
			prettyfmt.ErrorF("Failed to add \"%s\" as an asset.\nError: ", asset_name)
			generalerrors.HandleError(err)
		}

		prettyfmt.CPrintf("Successfully added \"%s\" as an asset.\nFormat: PDF\nSize: %d bytes\n", prettyfmt.GREEN, asset_name, len(file_data))
		return true
	default:
		prettyfmt.ErrorF("Failed to add \"%s\" as an asset.\nError: Unknown format!\n", asset_name)
		return false
	}
}

// This function will remove an asset from the wallet.
func (w *Wallet) Remove_Asset(asset_name string) bool {
	a := w.get_asset(asset_name)

	if a != nil {
		err := os.Remove(prettyfmt.SPathF(w.database_dir, asset_name))

		if err != nil {
			generalerrors.HandleError(err)
			return false
		}
	} else {
		return false
	}

	// This removes the asset from the slice of assets
	for i, as := range w.assets {
		if as == a {
			w.assets = append(w.assets[:i], w.assets[i+1:]...)
			break
		}
	}

	return true
}

func (w *Wallet) Get_Username() string {
	return w.username
}

func (w *Wallet) Confirm_Password(pass string) bool {
	pass_bytes := sha256.Sum256([]byte(pass))

	pass_bytes_str := prettyfmt.Sprintf("%X", pass_bytes)

	if w.password == pass_bytes_str {
		w.password = pass_bytes_str
		return true
	}

	return false
}
