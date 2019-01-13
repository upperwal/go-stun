package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/libp2p/go-libp2p-host"

	"github.com/libp2p/go-libp2p-swarm"

	cid "github.com/ipfs/go-cid"
	logging "github.com/ipfs/go-log"
	libp2p "github.com/libp2p/go-libp2p"
	crypto "github.com/libp2p/go-libp2p-crypto"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	inet "github.com/libp2p/go-libp2p-net"
	pstore "github.com/libp2p/go-libp2p-peerstore"
	ma "github.com/multiformats/go-multiaddr"
	mh "github.com/multiformats/go-multihash"
	"github.com/upperwal/go-libp2p-quic-transport" // packetConn branch
)

// IPFS bootstrap nodes. Used to find other peers in the network.
var bootstrapPeers = []string{
	/* "/ip4/104.131.131.82/tcp/4001/ipfs/QmaCpDMGvV2BGHeYERUEnRQAwe3N8SzbUtfsmvsqQLuvuJ",
	"/ip4/104.236.179.241/tcp/4001/ipfs/QmSoLPppuBtQSGwKDZT2M73ULpjvfd3aZ6ha4oFGL1KrGM",
	"/ip4/104.236.76.40/tcp/4001/ipfs/QmSoLV4Bbm51jM9C4gDYZQ9Cy3U6aXMJDAbzgu2fzaDs64",
	"/ip4/128.199.219.111/tcp/4001/ipfs/QmSoLSafTMBsPKadTEgaXctDQVcqN88CNLHXMkTNwMKPnu",
	"/ip4/178.62.158.247/tcp/4001/ipfs/QmSoLer265NRgSp2LA3dPaeykiS1J6DifTC88f5uVQKNAd", */
	"/ip4/104.43.165.178/udp/4000/quic/p2p/QmVbcMycaK8ni5CeiM7JRjBRAdmwky6dQ6KcoxLesZDPk9", // might be offline
}

var stunServer = "/ip4/104.43.165.178/udp/3000/" // might be offline

var rendezvous = "meet me hesfdfe"

func handleStream(stream inet.Stream) {
	log.Println("Got a new stream!")

	// Create a buffer stream for non blocking read and write.
	rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

	go readData(rw)
	go writeData(rw)

	// 'stream' will stay open until you close it (or the other side closes it).
}
func readData(rw *bufio.ReadWriter) {
	for {
		str, _ := rw.ReadString('\n')

		if str == "" {
			return
		}
		if str != "\n" {
			// Green console colour: 	\x1b[32m
			// Reset console colour: 	\x1b[0m
			fmt.Printf("\x1b[32m%s\x1b[0m> ", str)
		}

	}
}

func writeData(rw *bufio.ReadWriter) {
	stdReader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		sendData, err := stdReader.ReadString('\n')

		if err != nil {
			panic(err)
		}

		rw.WriteString(fmt.Sprintf("%s\n", sendData))
		rw.Flush()
	}
}

func main() {
	logging.SetLogLevel("dht", "DEBUG")
	logging.SetLogLevel("stun", "DEBUG")
	logging.SetLogLevel("relay", "DEBUG")
	logging.SetLogLevel("swarm2", "DEBUG")
	help := flag.Bool("h", false, "Display Help")
	rendezvousString := flag.String("r", rendezvous, "Unique string to identify group of nodes. Share this with your friends to let them connect with you")
	flag.Parse()

	if *help {
		fmt.Printf("This program demonstrates a simple p2p chat application using libp2p\n\n")
		fmt.Printf("Usage: Run './chat in two different terminals. Let them connect to the bootstrap nodes, announce themselves and connect to the peers\n")

		os.Exit(0)
	}

	ctx := context.Background()

	// libp2p.New constructs a new libp2p Host.
	// Other options can be added here.

	prvKey, _, err := crypto.GenerateKeyPair(crypto.RSA, 2048)
	if err != nil {
		panic(err)
	}

	stunMA, err := ma.NewMultiaddr(stunServer)
	if err != nil {
		panic(err)
	}

	quicOption := libp2pquic.TransportOpt{
		EnableStun:  true,
		StunServers: []ma.Multiaddr{stunMA},
	}

	quicTransport, err := libp2pquic.NewTransport(prvKey, quicOption)
	if err != nil {
		panic(err)
	}

	/* pc, err := libp2pquic.GetConnForAddr(quicTransport, "udp4", "0.0.0.0:0")
	if err != nil {
		panic(err)
	} */

	/* sclient, err := stun.NewClient(prvKey, pc)
	if err != nil {
		panic(err)
	} */

	/* stunMA, err := ma.NewMultiaddr(stunServer)
	if err != nil {
		panic(err)
	} */
	/* sclient.ConnectSTUNServer([]ma.Multiaddr{stunMA}) */

	host, err := libp2p.New(
		ctx,
		libp2p.ChainOptions(
			libp2p.Transport(quicTransport)),
		/* libp2p.Transport(tcp.NewTCPTransport)), */
		libp2p.ListenAddrStrings("/ip4/0.0.0.0/udp/0/quic"),
		libp2p.Identity(prvKey),
		//libp2p.EnableRelay(circuit.OptDiscovery),
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("Complete: ", host.ID().Pretty())
	fmt.Println("This node: ", host.ID())

	// Set a function as stream handler.
	// This function is called when a peer initiate a connection and starts a stream with this peer.
	host.SetStreamHandler("/chat/1.1.0", handleStream)

	kadDht, err := dht.New(ctx, host)
	if err != nil {
		panic(err)
	}

	// Let's connect to the bootstrap nodes first. They will tell us about the other nodes in the network.
	for _, peerAddr := range bootstrapPeers {
		addr, err := ma.NewMultiaddr(peerAddr)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(addr)
		peerinfo, _ := pstore.InfoFromP2pAddr(addr)

		if err := host.Connect(ctx, *peerinfo); err != nil {
			fmt.Println("Conn to bootstrap: ", err)
		} else {
			fmt.Println("Connection established with bootstrap node: ", *peerinfo)
		}
	}

	// We use a rendezvous point "meet me here" to announce our location.
	// This is like telling your friends to meet you at the Eiffel Tower.
	v1b := cid.V1Builder{Codec: cid.Raw, MhType: mh.SHA2_256}
	rendezvousPoint, _ := v1b.Sum([]byte(*rendezvousString))

	fmt.Println("announcing ourselves...")
	tctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	if err := kadDht.Provide(tctx, rendezvousPoint, true); err != nil {
		panic(err)
	}

	time.Sleep(5 * time.Second)

	// Now, look for others who have announced
	// This is like your friend telling you the location to meet you.
	// 'FindProviders' will return 'PeerInfo' of all the peers which
	// have 'Provide' or announced themselves previously.
	fmt.Println("searching for other peers...")
	tctx, cancel = context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	peers, err := kadDht.FindProviders(tctx, rendezvousPoint)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Found %d peers!\n", len(peers))

	for _, p := range peers {
		fmt.Println("Peer: ", p)
	}

	for _, p := range peers {
		if p.ID == host.ID() || len(p.Addrs) == 0 {
			// No sense connecting to ourselves or if addrs are not available
			continue
		}

		host.Peerstore().ClearAddrs(p.ID)

		var raddr ma.Multiaddr
		for _, addr := range p.Addrs {
			if strings.Contains(addr.String(), "p2p-circuit") || strings.Contains(addr.String(), "127.0.") || strings.Contains(addr.String(), "10.0") || strings.Contains(addr.String(), "127.0") || strings.Contains(addr.String(), "192.168") {
				continue
			} else {
				raddr = addr
			}
		}

		quicMA, _ := ma.NewMultiaddr("/quic")
		host.Peerstore().AddAddr(p.ID, raddr, pstore.PermanentAddrTTL)
		raddr = raddr.Decapsulate(quicMA)
		fmt.Println("Dialing to: ", p.ID, raddr)

		/* wait, err := sclient.PunchHole(raddr)
		if err != nil {
			panic(err)
		} */

		stream, err := host.NewStream(ctx, p.ID, "/chat/1.1.0")
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Println("Connected to: ", p)
		rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

		go writeData(rw)
		go readData(rw)

		/* var stream inet.Stream
		if <-wait {
			host.Network().(*swarm.Swarm).Backoff().Clear(p.ID)

			if err != nil {
				fmt.Println("direct connection failed", err)
				connectThroughRelay(ctx, host, p)
			} else {
				fmt.Println("Connected to: ", p)
				rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

				go writeData(rw)
				go readData(rw)
			}
		} else {
			fmt.Println("hole punching failed. connecting through relay")
			connectThroughRelay(ctx, host, p)
		} */

	}

	select {}
}

func connectThroughRelay(ctx context.Context, host host.Host, p pstore.PeerInfo) {
	host.Network().(*swarm.Swarm).Backoff().Clear(p.ID)
	relayaddr, err := ma.NewMultiaddr("/p2p-circuit/ipfs/" + p.ID.Pretty())
	if err != nil {
		panic(err)
	}

	pRelayInfo := pstore.PeerInfo{
		ID:    p.ID,
		Addrs: []ma.Multiaddr{relayaddr},
	}

	if err := host.Connect(context.Background(), pRelayInfo); err != nil {
		fmt.Println("relay connect failed", err)
	} else {
		stream, err := host.NewStream(ctx, p.ID, "/chat/1.1.0")

		if err != nil {
			fmt.Println("relay stream failed", err)
		} else {
			fmt.Println("Connected to: ", p)
			rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

			go writeData(rw)
			go readData(rw)
		}
	}
}
