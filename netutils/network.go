package netutils

import (
	"context"
	"errors"
	"github.com/TheJ0lly/GoChain/asset"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/multiformats/go-multiaddr"
	"log"
)

const DEFAULTALLIP4 = "/ip4/0.0.0.0/tcp/0"
const DEFAULTALLIP6 = "/ip6/::/tcp/0"

type Options struct {
	privKey crypto.PrivKey
	useIP6  bool
	useIP4  bool
	addrs   []string
}

func CreateNodeOptions(PrivateKey crypto.PrivKey, IPv4 bool, IPv6 bool, addresses ...string) Options {
	return Options{
		privKey: PrivateKey,
		useIP6:  IPv6,
		useIP4:  IPv4,
		addrs:   addresses,
	}
}

// getDefaultAddresses - this function will assign the multi address form to notify the libp2p that it should auto search for an address
func getDefaultAddresses(Opt Options) []string {
	var addrs []string

	if Opt.useIP4 {
		addrs = append(addrs, DEFAULTALLIP4)
	}

	if Opt.useIP6 {
		addrs = append(addrs, DEFAULTALLIP6)
	}

	return addrs
}

func CreateNewNode(Opt Options) (core.Host, error) {
	defAddrs := getDefaultAddresses(Opt)

	if len(Opt.addrs) == 0 {
		log.Printf("INFO: no addresses have been given - trying the options given...\n")

		if len(defAddrs) == 0 {
			log.Printf("INFO: auto searching for IPv4 and IPv6 has not been allowed\n")
			return nil, errors.New("cannot create new node")
		}

	} else {
		log.Printf("INFO: given addresses: %v\n", Opt.addrs)
		//If there are default addresses, we add them together with the defaults, if they are selected
		defAddrs = append(defAddrs, Opt.addrs...)
	}

	h, err := libp2p.New(
		libp2p.ListenAddrStrings(defAddrs...),
		libp2p.Identity(Opt.privKey),
	)

	if err != nil {
		return nil, err
	}

	return h, nil
}

// MakeRequest - this function will make the request to the known addresses for a specific asset.
func MakeRequest(addresses []multiaddr.Multiaddr, ha core.Host, assetName string) (bool, *asset.Asset) {
	if addresses == nil {
		log.Printf("INFO: no known addresses - aborting request for %s\n", assetName)
		return false, nil
	}

	ok := false
	var s network.Stream
	var buff []byte
	for _, addr := range addresses {
		info, err := peer.AddrInfoFromP2pAddr(addr)

		if err != nil {
			log.Printf("ERROR: %s\n", err)
			log.Printf("INFO: moving to the next address\n")
			continue
		}

		ha.Peerstore().AddAddrs(info.ID, info.Addrs, peerstore.AddressTTL)
		log.Printf("INFO: trying to connect to %s\n", addr.String())

		s, err = ha.NewStream(context.Background(), info.ID, "REQUEST")

		if err != nil {
			log.Printf("ERROR: %s\n", err)
			log.Printf("INFO: moving to the next address\n")
			continue
		}

		log.Printf("INFO: connected successfully to %s\n", addr.String())

		log.Printf("INFO: requesting asset - %s\n", assetName)

		_, err = s.Write([]byte(assetName))

		if err != nil {
			log.Printf("ERROR: %s\n", err)
			log.Printf("INFO: moving to the next address\n")
			continue
		}

		resp := make([]byte, 10)
		log.Printf("INFO: waiting for response code\n")
		_, err = s.Read(resp)

		if err != nil {
			log.Printf("ERROR: %s\n", err)
			continue
		}
		val := GetNumberFromResponse(resp)
		log.Printf("INFO: received - %d\n", val)

		if val == FailedConversion {
			log.Printf("INFO: moving to the next address\n")
			continue
		} else if val == AssetNotFound {
			log.Printf("INFO: current node does not have asset %s\n", assetName)
			log.Printf("INFO: moving to the next address\n")
			continue
		} else {
			ok, buff = ReceiveAsset(s, assetName, val)
			if !ok {
				log.Printf("INFO: moving to the next address\n")
				continue
			}
		}

		log.Printf("INFO: request executed successfully\n")
		ok = true
		break
	}

	log.Printf("INFO: closing the network stream\n")
	err := s.Close()

	if err != nil {
		log.Printf("ERROR: %s\n", err)
		ok = false
	}

	if ok {
		ft := asset.DetermineType(buff)

		if ft == asset.UNKNOWN {
			ok = false
		}

		as := asset.CreateNewAsset(assetName, asset.DetermineType(buff), buff)
		return true, as
	} else {
		return false, nil
	}
}
