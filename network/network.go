package network

import (
	"context"
	"errors"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/multiformats/go-multiaddr"
	"log"
)

const DEFAULTALLIP4 = "/ip4/0.0.0.0/tcp/0"
const DEFAULTALLIP6 = "/ip6/::/tcp/0"

type Options struct {
	useIP6 bool
	useIP4 bool
	addrs  []string
}

func CreateNodeOptions(IPv4 bool, IPv6 bool, addresses ...string) Options {
	return Options{
		useIP6: IPv6,
		useIP4: IPv4,
		addrs:  addresses,
	}
}

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
		log.Printf("No addresses have been given, trying the options given...\n")

		if len(defAddrs) == 0 {
			log.Printf("Auto searching for IPv4 and IPv6 has not been allowed\n")
			return nil, errors.New("cannot create new node")
		}

	} else {
		log.Printf("Given addresses: %v\n", Opt.addrs)
		//If there are default addresses, we add them together with the defaults, if they are selected
		defAddrs = append(defAddrs, Opt.addrs...)
	}

	h, err := libp2p.New(
		libp2p.ListenAddrStrings(defAddrs...),
	)

	if err != nil {
		return nil, err
	}

	log.Printf("Initializing node with the following addresses: %v\n", h.Addrs())

	log.Printf("Node has been initialized!\n")
	return h, nil
}

func SendTo(from host.Host, data []byte, to multiaddr.Multiaddr) error {
	toInfo, err := peer.AddrInfoFromP2pAddr(to)

	if err != nil {
		return err
	}
	from.Peerstore().AddAddrs(toInfo.ID, toInfo.Addrs, peerstore.TempAddrTTL)
	s, err := from.NewStream(context.Background(), toInfo.ID, "TRANSFER")

	if err != nil {
		return err
	}

	log.Printf("Sending data...\n")
	//Write the data len to the stream, to help the "other side" to make enough space for this data.

	//Write the data to the stream
	n, err := s.Write(data)

	if err != nil {
		return err
	}

	log.Printf("Sent %d bytes over to: %v\n", n, to)

	return nil
}
