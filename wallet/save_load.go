package wallet

import (
	"crypto/rsa"
	"crypto/sha256"
	"encoding/json"
	"os"

	"github.com/TheJ0lly/GoChain/asset"
	"github.com/TheJ0lly/GoChain/generalerrors"
	"github.com/TheJ0lly/GoChain/prettyfmt"
)

type Wallet_IE struct {
	Username     string         `json:"USERNAME"`
	Password     string         `json:"PASSWORD"`
	Pub_Key      rsa.PublicKey  `json:"PUB_KEY"`
	Priv_Key     rsa.PrivateKey `json:"PRIV_KEY"`
	Database_Dir string         `json:"DB_DIR"`
	Assets       []string       `json:"ASSETS"`
}

func (w *Wallet) Save_State() error {

	pass_bytes := sha256.Sum256([]byte(w.password))

	pass_bytes_str := prettyfmt.Sprintf("%X", pass_bytes)

	wie := &Wallet_IE{Username: w.username, Password: pass_bytes_str, Pub_Key: w.public_key, Priv_Key: w.private_key, Database_Dir: w.database_dir}

	for _, a := range w.assets {
		wie.Assets = append(wie.Assets, a.Get_Name())
	}

	bytes_to_write, err := json.Marshal(wie)

	if err != nil {
		return &generalerrors.JSONMarshalFailed{Object: "Wallet"}
	}

	err = os.WriteFile("./ws", bytes_to_write, 0666)

	if err != nil {
		return &generalerrors.WriteFileFailed{File: "./ws"}
	}

	return nil
}

func Load_Wallet() (*Wallet, error) {
	prettyfmt.CPrint("Looking for save of the wallet...\n", prettyfmt.BLUE)

	bytes_read, err := os.ReadFile("./ws")

	if err != nil {
		return nil, &generalerrors.ReadFileFailed{File: "./ws"}
	}

	wie := &Wallet_IE{}

	err = json.Unmarshal(bytes_read, wie)

	if err != nil {
		return nil, &generalerrors.JSONUnMarshalFailed{Object: "Wallet"}
	}

	files, err := os.ReadDir(wie.Database_Dir)

	if err != nil {
		return nil, &generalerrors.ReadDirFailed{Dir: wie.Database_Dir}
	}

	asset_len := len(wie.Assets)

	assets_to_recreate := []string{}

	if asset_len != len(files) {
		prettyfmt.CPrint("Some assets have been corrupted or deleted! Recreating the assets that can be found...\n", prettyfmt.YELLOW)

		assets_to_recreate = append(assets_to_recreate, wie.Assets...)
	}

	w := &Wallet{username: wie.Username, password: wie.Password, public_key: wie.Pub_Key, private_key: wie.Priv_Key, database_dir: wie.Database_Dir}

	for _, f := range files {
		continue_to_recreate := false

		if len(assets_to_recreate) == 0 {
			continue_to_recreate = true
		} else {
			for _, a := range assets_to_recreate {
				if f.Name() == a {
					continue_to_recreate = true
					break
				}
			}
		}

		if !continue_to_recreate {
			file_to_remove := prettyfmt.Sprintf("%s/%s", wie.Database_Dir, f.Name())
			err = os.Remove(file_to_remove)

			if err != nil {
				return nil, &generalerrors.RemoveFileFailed{File: file_to_remove}
			}

			continue
		}

		file_to_recreate := prettyfmt.Sprintf("%s/%s", wie.Database_Dir, f.Name())

		bytes_read, err = os.ReadFile(file_to_recreate)

		if err != nil {
			return nil, &generalerrors.ReadFileFailed{File: file_to_recreate}
		}

		ft := asset.Determine_Asset_Type(bytes_read)

		if ft == asset.UNKNOWN {
			prettyfmt.CPrintf("This asset may have been corrupted, changed, or added manually! Skipping - %s\n", prettyfmt.YELLOW, f.Name())
			continue
		}

		asset := asset.Create_New_Asset(f.Name(), ft, bytes_read)

		w.assets = append(w.assets, asset)
	}

	return w, nil
}
