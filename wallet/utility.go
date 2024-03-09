package wallet

import (
	"github.com/TheJ0lly/GoChain/asset"
	"github.com/multiformats/go-multiaddr"
	"strings"
)

// This function check if the asset exists based on its name.
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

func (w *Wallet) checkIfAddrExists(ma multiaddr.Multiaddr) bool {
	for _, ad := range w.mKnownHosts {
		if ad.String() == ma.String() {
			return true
		}
	}
	return false
}

func checkAddrIsNotLH(ma multiaddr.Multiaddr) bool {
	var IP string

	maStr := ma.String()

	for i := 5; i < len(maStr); i++ {
		if maStr[i] == '/' {
			break
		}
		IP += string(maStr[i])
	}

	if strings.Contains(IP, "127.0") || strings.Contains(IP, "::1") {
		return false
	}

	return true
}
