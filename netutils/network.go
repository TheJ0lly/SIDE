package netutils

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/TheJ0lly/GoChain/asset"
	"github.com/TheJ0lly/GoChain/blockchain"
	"github.com/TheJ0lly/GoChain/metadata"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/multiformats/go-multiaddr"
	"io"
	"log"
	"net"
	"strconv"
	"time"
)

type Options struct {
	privKey crypto.PrivKey
	ip      string
	port    string
}

func CreateNodeOptions(PrivateKey crypto.PrivKey, IP string, Port string) Options {
	return Options{
		privKey: PrivateKey,
		ip:      IP,
		port:    Port,
	}
}

func CreateNewNode(Opt Options) (core.Host, error) {
	var ipversion string
	var addresses []string
	var port string

	if Opt.port == "" {
		log.Printf("INFO: no port has been given - auto searching for a port\n")
		port = "0"
	} else {
		log.Printf("INFO: given port %s\n", Opt.port)
		port = Opt.port
	}

	if Opt.ip != "" {
		log.Printf("INFO: parsing given IP\n")
		IP := net.ParseIP(Opt.ip)

		if IP == nil {
			return nil, errors.New("could not parse IP given")
		}

		if IP.To4() != nil {
			ipversion = "ip4"
		} else {
			ipversion = "ip6"
		}

		addresses = append(addresses, fmt.Sprintf("/%s/%s/tcp/%s", ipversion, Opt.ip, Opt.port))
	} else {
		log.Printf("INFO: no IP has been given - auto searching for IPv4 address\n")

		addresses = append(addresses, fmt.Sprintf("/ip4/0.0.0.0/tcp/%s", port))
		addresses = append(addresses, fmt.Sprintf("/ip6/::/tcp/%s", port))
	}

	h, err := libp2p.New(
		libp2p.ListenAddrStrings(addresses...),
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
	var end time.Duration
	for _, addr := range addresses {
		start := time.Now()
		log.Printf("INFO: getting peer information\n")
		info, err := peer.AddrInfoFromP2pAddr(addr)

		if err != nil {
			log.Printf("ERROR: %s\n", err)
			log.Printf("INFO: moving to the next address\n")
			continue
		}

		ha.Peerstore().AddAddrs(info.ID, info.Addrs, peerstore.TempAddrTTL)
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
			ok, buff = receiveAsset(s, assetName, val)
			end = time.Since(start)

			if !ok {
				log.Printf("INFO: moving to the next address\n")
				continue
			}
		}

		log.Printf("INFO: request executed successfully - %v ms\n", end.Milliseconds())
		ok = true
		break
	}

	if s != nil {
		log.Printf("INFO: closing the network stream\n")
		err := s.Close()

		if err != nil {
			log.Printf("ERROR: %s\n", err)
			ok = false
		}
	}

	if ok {
		as := asset.CreateNewAsset(assetName, buff)
		return true, as
	} else {
		return false, nil
	}
}

func CreateNewBlockchainFromConn(ha core.Host, dbLoc string, ma multiaddr.Multiaddr) (*blockchain.BlockChain, error) {
	log.Printf("INFO: getting peer information\n")
	info, err := peer.AddrInfoFromP2pAddr(ma)

	if err != nil {
		return nil, err
	}

	ha.Peerstore().AddAddrs(info.ID, info.Addrs, peerstore.TempAddrTTL)
	log.Printf("INFO: trying to connect to %s\n", ma.String())

	s, err := ha.NewStream(context.Background(), info.ID, "INITIALIZE")

	if err != nil {
		return nil, err
	}

	log.Printf("INFO: connected successfully to %s\n", ma.String())

	log.Printf("INFO: requesting number of blocks\n")
	resp := make([]byte, 10)

	_, err = s.Read(resp)

	if err != nil {
		return nil, err
	}

	val := GetNumberFromResponse(resp)
	log.Printf("INFO: received - %d\n", val)

	if val == FailedConversion {
		return nil, errors.New("failed to convert response to int\n")
	}

	blocks := InitializeProtocol(val, s)

	if blocks == nil {
		return nil, errors.New("there are no blocks")
	}

	bc := blockchain.CreateNewBlockchainFromData(dbLoc, blocks)

	return bc, nil
}

func InitializeProtocol(numBlocks int, s network.Stream) []*blockchain.Block {
	var bcc = make([]*blockchain.Block, 0)
	log.Printf("INFO: waiting to read %d blocks\n", numBlocks)
	for i := 0; i < numBlocks; i++ {
		resp := make([]byte, 10)

		log.Printf("INFO: waiting for block size\n")
		_, err := s.Read(resp)

		if err != nil {
			log.Printf("ERROR: %s\n", err)
			return nil
		}
		l := GetNumberFromResponse(resp)

		if l == FailedConversion {
			log.Printf("ERROR: failed to convert response to int\n")
			return nil
		}
		log.Printf("INFO: received - %d\n", l)

		log.Printf("INFO: creating buffer of capacity - %d\n", l)
		stor := make([]byte, l)

		log.Printf("INFO: attempting to read %d bytes\n", l)
		n, err := io.ReadFull(bufio.NewReader(s), stor)

		if err != nil {
			log.Printf("ERROR: %s\n", err)
			return nil
		}

		log.Printf("INFO: read %d bytes from the connection\n", n)
		if n != l {
			log.Printf("ERROR: read a different amount of bytes than expected\n")
			return nil
		}

		b := blockchain.ImportBlockFromConn(stor)

		log.Printf("INFO: received block - %X\n", b.GetBlockTreeMatrix().RootHash)

		if b == nil {
			return nil
		}

		bcc = append(bcc, b)
	}

	return bcc
}

func FloodProtocol(addresses []multiaddr.Multiaddr, h core.Host, md *metadata.MetaData) {
	if addresses == nil {
		log.Printf("INFO: no known addresses - aborting flood\n")
		return
	}

	log.Printf("INFO: starting the flood protocol\n")

	for _, addr := range addresses {
		log.Printf("INFO: getting peer information\n")
		info, err := peer.AddrInfoFromP2pAddr(addr)

		if err != nil {
			log.Printf("ERROR: %s\n", err)
			log.Printf("INFO: moving to the next address\n")
			continue
		}

		h.Peerstore().AddAddrs(info.ID, info.Addrs, peerstore.TempAddrTTL)
		log.Printf("INFO: trying to connect to %s\n", addr.String())

		s, err := h.NewStream(context.Background(), info.ID, "FLOOD")

		if err != nil {
			log.Printf("ERROR: %s\n", err)
			log.Printf("INFO: moving to the next address\n")
			continue
		}

		log.Printf("INFO: connected successfully to %s\n", addr.String())

		mie := metadata.MetadataIE{
			Source:      md.GetSourceName(),
			Destination: md.GetDestinationName(),
			AssetName:   md.GetAssetName(),
			Hash:        md.GetMetadataHash(),
		}

		b, err := json.Marshal(mie)

		if err != nil {
			log.Printf("ERROR: %s\n", err)
			log.Printf("INFO: moving to the next address\n")
			_, err = s.Write([]byte("-1"))

			if err != nil {
				log.Printf("ERROR: %s\n", err)
			}
			continue
		}

		log.Printf("INFO: sending length of metadata\n")
		_, err = s.Write([]byte(strconv.Itoa(len(b))))

		if err != nil {
			log.Printf("ERROR: %s\n", err)
			log.Printf("INFO: moving to the next address\n")
			continue
		}

		log.Printf("INFO: waiting for ready signal\n")
		var stor = make([]byte, 10)

		_, err = s.Read(stor)

		if err != nil {
			log.Printf("ERROR: %s\n", err)
			log.Printf("INFO: moving to the next address\n")
			continue
		}

		resp := ConvertBytesToString(stor)

		if resp == "READY" {
			log.Printf("INFO: ready signal recieved - sending bytes\n")
			_, err = s.Write(b)

			if err != nil {
				log.Printf("ERROR: %s\n", err)
				log.Printf("INFO: moving to the next address\n")
				continue
			}
		}

		log.Printf("INFO: updated node %s\n", addr.String())
		log.Printf("INFO: moving to the next known node\n")

	}
}

func ForwardMetadata(addresses []multiaddr.Multiaddr, h core.Host, metadata []byte, jumpAddr string) {
	jumpA, err := multiaddr.NewMultiaddr(jumpAddr)

	if err != nil {
		log.Printf("ERROR: %s\n", err)
		return
	}

	for _, addr := range addresses {
		if addr.Equal(jumpA) {
			log.Printf("INFO: jumping sender address\n")
			continue
		}
		log.Printf("INFO: getting peer information\n")
		info, err := peer.AddrInfoFromP2pAddr(addr)

		if err != nil {
			log.Printf("ERROR: %s\n", err)
			log.Printf("INFO: moving to the next address\n")
			continue
		}

		h.Peerstore().AddAddrs(info.ID, info.Addrs, peerstore.TempAddrTTL)
		log.Printf("INFO: trying to connect to %s\n", addr.String())

		s, err := h.NewStream(context.Background(), info.ID, "FLOOD")

		if err != nil {
			log.Printf("ERROR: %s\n", err)
			log.Printf("INFO: moving to the next address\n")
			continue
		}

		log.Printf("INFO: connected successfully to %s\n", addr.String())

		_, err = s.Write([]byte(strconv.Itoa(len(metadata))))

		if err != nil {
			log.Printf("ERROR: %s\n", err)
			log.Printf("INFO: moving to the next address\n")
			continue
		}

		log.Printf("INFO: waiting for ready signal\n")
		var stor = make([]byte, 10)

		_, err = s.Read(stor)

		if err != nil {
			log.Printf("ERROR: %s\n", err)
			log.Printf("INFO: moving to the next address\n")
			continue
		}

		resp := ConvertBytesToString(stor)

		if resp == "READY" {
			log.Printf("INFO: ready signal recieved - sending bytes\n")
			_, err = s.Write(metadata)

			if err != nil {
				log.Printf("ERROR: %s\n", err)
				log.Printf("INFO: moving to the next address\n")
				continue
			}
		}

		log.Printf("INFO: successfully forward metadata to %s\n", addr.String())
	}

	log.Printf("INFO: forwarded metadata to all known nodes\n")
}

func GetHostAddressFromConnection(conn core.Conn) string {
	hostAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/p2p/%s", conn.RemotePeer()))

	return conn.RemoteMultiaddr().Encapsulate(hostAddr).String()
}
