package wallet

import (
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/TheJ0lly/GoChain/asset"
	"github.com/TheJ0lly/GoChain/generalerrors"
	"github.com/TheJ0lly/GoChain/netutils"
	"github.com/TheJ0lly/GoChain/osspecifics"
	"github.com/libp2p/go-libp2p/core"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/multiformats/go-multiaddr"
	"log"
	"os"
)

const (
	//64 bits (Source name) + 8 bits (Space) + 64 bits (Destination name) + 8 bits (Space) + 64 bits (Asset Name) + 8 bits (Empty/NULL)
	usernameMaxLength = 8
)

type Wallet struct {
	mUsername    string
	mPassword    [32]byte
	mPublicKey   crypto.PubKey
	mPrivateKey  crypto.PrivKey
	mDatabaseDir string
	mAssets      []*asset.Asset
	mHost        core.Host
	mKnownHosts  []core.Multiaddr
}

// CreateNewWallet - This function will create a wallet.
func CreateNewWallet(username string, password string, dbLoc string, IP string, Port string) (*Wallet, error) {
	if len(username) > usernameMaxLength {
		return nil, &generalerrors.UsernameTooLong{Length: usernameMaxLength}
	}

	privateKey, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, rand.Reader)

	if err != nil {
		return nil, err
	}

	host, err := netutils.CreateNewNode(netutils.CreateNodeOptions(privateKey, IP, Port))

	if err != nil {
		return nil, err
	}

	passBytes := sha256.Sum256([]byte(password))
	w := &Wallet{mUsername: username, mPassword: passBytes, mPrivateKey: privateKey, mPublicKey: privateKey.GetPublic(), mAssets: nil, mDatabaseDir: dbLoc, mHost: host, mKnownHosts: nil}
	log.Printf("INFO: node has been initialized with address: %v\n", w.GetHostAddress())

	return w, nil
}

// AddAssetFromLocal - This function will add an asset to the wallet.
func (w *Wallet) AddAssetFromLocal(assetName string, fileLocation string) (*asset.Asset, error) {

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

	var assetToAdd *asset.Asset = asset.CreateNewAsset(assetName, fileData)
	AssetPath := osspecifics.CreatePath(w.mDatabaseDir, assetName)

	err = os.WriteFile(AssetPath, fileData, 0444)

	if err != nil {
		return nil, &generalerrors.WriteFileFailed{File: AssetPath}
	}

	w.mAssets = append(w.mAssets, assetToAdd)
	return assetToAdd, nil
}

func (w *Wallet) AddAssetFromNode(as *asset.Asset) (*asset.Asset, error) {
	AssetPath := osspecifics.CreatePath(w.mDatabaseDir, as.GetName())
	err := os.WriteFile(AssetPath, as.GetAssetBytes(), 0444)

	if err != nil {
		return nil, &generalerrors.WriteFileFailed{File: AssetPath}
	}
	w.mAssets = append(w.mAssets, as)
	return as, nil
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

func (w *Wallet) AddNode(address string) (multiaddr.Multiaddr, error) {
	ma, err := multiaddr.NewMultiaddr(address)

	if err != nil {
		return nil, err
	}

	if w.checkIfAddrExists(ma) {
		return nil, errors.New(fmt.Sprintf("the address %s is already added", ma.String()))
	}

	w.mKnownHosts = append(w.mKnownHosts, ma)
	return ma, nil
}

func (w *Wallet) GetAsset(assetName string) (*asset.Asset, error) {
	for _, a := range w.mAssets {
		if a.GetName() == assetName {
			return a, nil
		}
	}

	return nil, &generalerrors.AssetDoesNotExist{AssetName: assetName}
}

func (w *Wallet) GetUsername() string {
	return w.mUsername
}

func (w *Wallet) ConfirmPassword(pass string) bool {
	passBytes := sha256.Sum256([]byte(pass))

	return w.mPassword == passBytes
}

func (w *Wallet) GetNodesAddresses() []multiaddr.Multiaddr {
	return w.mKnownHosts
}

func (w *Wallet) GetDBLocation() string {
	return w.mDatabaseDir
}

func (w *Wallet) GetHost() core.Host {
	return w.mHost
}

func (w *Wallet) GetPrivateKey() crypto.PrivKey { return w.mPrivateKey }

func (w *Wallet) GetHostAddress() string {
	hostAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/p2p/%s", w.mHost.ID()))
	addrs := w.mHost.Addrs()

	var addr multiaddr.Multiaddr

	for _, a := range addrs {
		if checkAddrIsNotLH(a) {
			addr = a
			break
		}
	}

	return addr.Encapsulate(hostAddr).String()

}
