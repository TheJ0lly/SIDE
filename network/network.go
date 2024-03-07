package network

import (
	"errors"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core"
	"github.com/libp2p/go-libp2p/core/crypto"
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

	log.Printf("INFO: node has been initialized with addresses: %v\n", h.Addrs())
	return h, nil
}
