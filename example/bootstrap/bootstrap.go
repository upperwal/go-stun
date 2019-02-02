package main

import (
	"context"
	"fmt"
	mrand "math/rand"

	logging "github.com/ipfs/go-log"
	floodsub "github.com/libp2p/go-floodsub"
	libp2p "github.com/libp2p/go-libp2p"
	circuit "github.com/libp2p/go-libp2p-circuit"
	"github.com/libp2p/go-libp2p-crypto"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	inet "github.com/libp2p/go-libp2p-net"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/upperwal/go-libp2p-quic-transport"
)

type netNotifiee struct{}

func (nn *netNotifiee) Connected(n inet.Network, c inet.Conn) {
	fmt.Printf("Connected to: %s/p2p/%s\n", c.RemoteMultiaddr(), c.RemotePeer().Pretty())
}

func (nn *netNotifiee) Disconnected(n inet.Network, v inet.Conn)   {}
func (nn *netNotifiee) OpenedStream(n inet.Network, v inet.Stream) {}
func (nn *netNotifiee) ClosedStream(n inet.Network, v inet.Stream) {}
func (nn *netNotifiee) Listen(n inet.Network, a ma.Multiaddr)      {}
func (nn *netNotifiee) ListenClose(n inet.Network, a ma.Multiaddr) {}

func main() {
	logging.SetLogLevel("dht", "DEBUG")
	logging.SetLogLevel("relay", "DEBUG")
	logging.SetLogLevel("pubsub", "DEBUG")
	ctx := context.Background()

	// libp2p.New constructs a new libp2p Host.
	// Other options can be added here.
	sourceMultiAddr, _ := ma.NewMultiaddr("/ip4/0.0.0.0/udp/4000/quic")

	r := mrand.New(mrand.NewSource(int64(10)))
	prvKey, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	if err != nil {
		panic(err)
	}
	tpt, err := libp2pquic.NewTransport(prvKey, libp2pquic.TransportOpt{})
	if err != nil {
		panic(err)
	}
	host, err := libp2p.New(
		ctx,
		libp2p.ListenAddrs(sourceMultiAddr),
		libp2p.Identity(prvKey),
		libp2p.Transport(tpt),
		libp2p.EnableRelay(circuit.OptHop),
	)
	if err != nil {
		panic(err)
	}

	no := &netNotifiee{}
	host.Network().Notify(no)

	fmt.Println("This node: ", host.ID().Pretty(), " ", host.Addrs())

	_, err = dht.New(ctx, host)
	if err != nil {
		panic(err)
	}

	_, err = floodsub.NewGossipSub(ctx, host)
	if err != nil {
		panic(err)
	}

	select {}
}
