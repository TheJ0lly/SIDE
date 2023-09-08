package wallet

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"os"

	"github.com/TheJ0lly/GoChain/asset"
)

type Wallet_IE struct {
	Username     string         `json:"USERNAME"`
	Password     string         `json:"PASSWORD"`
	Pub_Key      rsa.PublicKey  `json:"PUB_KEY"`
	Priv_Key     rsa.PrivateKey `json:"PRIV_KEY"`
	Database_Dir string         `json:"DB_DIR"`
	Assets       []string       `json:"ASSETS"`
}

func (w *Wallet) Save_State() {
	wie := &Wallet_IE{Username: w.username, Password: w.password, Pub_Key: w.public_key, Priv_Key: w.private_key, Database_Dir: w.database_dir}

	for _, a := range w.assets {
		wie.Assets = append(wie.Assets, a.Get_Name())
	}

	bytes_to_write, err := json.Marshal(wie)

	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return
	}

	err = os.WriteFile("./ws", bytes_to_write, 0666)

	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return
	}
}

func load_wallet() *Wallet {
	bytes_read, err := os.ReadFile("./ws")

	if err != nil {
		fmt.Printf("There is no save file! Recreating wallet...\n")
		return nil
	}

	wie := &Wallet_IE{}

	err = json.Unmarshal(bytes_read, wie)

	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return nil
	}

	files, err := os.ReadDir(wie.Database_Dir)

	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return nil
	}

	asset_len := len(wie.Assets)

	assets_to_recreate := []string{}

	if asset_len != len(files) {
		fmt.Printf("Some assets have been corrupted or deleted! Recreating the assets that can be found...\n")

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
			os.Remove(fmt.Sprintf("%s\\%s", wie.Database_Dir, f.Name()))
			continue
		}

		bytes_read, err = os.ReadFile(fmt.Sprintf("%s\\%s", wie.Database_Dir, f.Name()))

		if err != nil {
			fmt.Printf("%s\n", err.Error())
			return nil
		}

		ft := asset.Determine_Asset_Type(bytes_read)

		if ft == asset.UNKNOWN {
			fmt.Printf("This asset may have been corrupted, changed, or added manually! Skipping - %s\n", f.Name())
		}

		asset := asset.Create_New_Asset(f.Name(), ft, bytes_read)

		w.assets = append(w.assets, asset)
	}

	return w
}
