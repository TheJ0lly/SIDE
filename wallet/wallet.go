package wallet

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"github.com/TheJ0lly/GoChain/osspecifics"
	"os"
	"path/filepath"

	"github.com/TheJ0lly/GoChain/asset"
	"github.com/TheJ0lly/GoChain/generalerrors"
)

const (
	//64 bits (Source name) + 8 bits (Space) + 64 bits (Destination name) + 8 bits (Space) + 64 bits (Asset Name) + 8 bits (Empty/NULL)
	keyBitSize        = 216
	usernameMaxLength = 8
)

type Wallet struct {
	mUsername    string
	mPassword    [32]byte
	mPublicKey   rsa.PublicKey
	mPrivateKey  rsa.PrivateKey
	mDatabaseDir string
	mAssets      []*asset.Asset
}

// CreateNewWallet - This function will create a wallet.
func CreateNewWallet(username string, password string, dbLoc string) (*Wallet, error) {
	if len(username) > usernameMaxLength {
		return nil, &generalerrors.UsernameTooLong{Length: usernameMaxLength}
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, keyBitSize)

	if err != nil {
		return nil, err
	}

	passBytes := sha256.Sum256([]byte(password))

	w := &Wallet{mUsername: username, mPassword: passBytes, mPrivateKey: *privateKey, mPublicKey: privateKey.PublicKey, mAssets: nil, mDatabaseDir: dbLoc}

	return w, nil
}

// AddAsset - This function will add an asset to the wallet.
func (w *Wallet) AddAsset(assetName string, fileLocation string) (*asset.Asset, error) {

	if w.checkAssetExists(assetName) {
		return nil, &generalerrors.AssetAlreadyExists{AssetName: assetName}
	}

	temp, err := os.Stat(fileLocation)

	if os.IsNotExist(err) {
		return nil, &generalerrors.AssetInitialLocationDoesNotExist{Location: fileLocation}
	}

	if temp.IsDir() {
		return nil, &generalerrors.AssetIsAFolder{Location: fileLocation}
	}

	fileData, err := os.ReadFile(fileLocation)

	if err != nil {
		return nil, &generalerrors.ReadFileFailed{File: fileLocation}
	}

	var assetToAdd *asset.Asset
	AssetPath := osspecifics.CreatePath(w.mDatabaseDir, assetName)

	switch asset.DetermineAssetType(fileData) {
	case asset.JPEG:
		assetToAdd = asset.CreateNewAsset(assetName, asset.JPEG, fileData)
	case asset.PDF:
		assetToAdd = asset.CreateNewAsset(assetName, asset.PDF, fileData)
	default:
		fileExt := filepath.Ext(AssetPath)
		return nil, &generalerrors.UnknownFormat{FileExt: fileExt}
	}

	err = os.WriteFile(AssetPath, fileData, 0444)

	if err != nil {
		return nil, &generalerrors.WriteFileFailed{File: AssetPath}
	}

	w.mAssets = append(w.mAssets, assetToAdd)
	return assetToAdd, nil
}

// RemoveAsset - This function will remove an asset from the wallet.
func (w *Wallet) RemoveAsset(assetName string) (*asset.Asset, error) {
	a := w.getAsset(assetName)

	if a != nil {
		path := osspecifics.CreatePath(w.mDatabaseDir, assetName)
		err := os.Remove(path)

		if err != nil {
			return nil, &generalerrors.RemoveFileFailed{File: path}
		}
	} else {
		return nil, &generalerrors.AssetDoesNotExist{AssetName: assetName}
	}

	var assetToRet *asset.Asset

	// This removes the asset from the slice of mAssets
	for i, as := range w.mAssets {
		if as == a {
			assetToRet = a
			w.mAssets = append(w.mAssets[:i], w.mAssets[i+1:]...)
			break
		}
	}

	return assetToRet, nil
}

func (w *Wallet) ViewAssets() []asset.Asset {
	var assetSlice []asset.Asset

	for _, a := range w.mAssets {
		assetSlice = append(assetSlice, a.GetAssetCopy())
	}

	return assetSlice
}

func (w *Wallet) GetUsername() string {
	return w.mUsername
}

func (w *Wallet) ConfirmPassword(pass string) bool {
	passBytes := sha256.Sum256([]byte(pass))

	return w.mPassword == passBytes
}

func (w *Wallet) GetDBLocation() string {
	return w.mDatabaseDir
}
