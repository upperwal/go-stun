package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/lucas-clemente/quic-go"
	"github.com/multiformats/go-multiaddr-net"

	logging "github.com/ipfs/go-log"
	crypto "github.com/libp2p/go-libp2p-crypto"
	ma "github.com/multiformats/go-multiaddr"
	stun "github.com/upperwal/go-stun"
)

var quicConfig = &quic.Config{
	Versions:                              []quic.VersionNumber{quic.VersionMilestone0_10_0},
	MaxIncomingStreams:                    1000,
	MaxIncomingUniStreams:                 -1,              // disable unidirectional streams
	MaxReceiveStreamFlowControlWindow:     3 * (1 << 20),   // 3 MB
	MaxReceiveConnectionFlowControlWindow: 4.5 * (1 << 20), // 4.5 MB
	AcceptCookie: func(clientAddr net.Addr, cookie *quic.Cookie) bool {
		// TODO(#6): require source address validation when under load
		return true
	},
	KeepAlive:   true,
	IdleTimeout: 30 * time.Hour,
}

func main() {

	logging.SetLogLevel("stun", "DEBUG")

	port := flag.String("p", "0", "listener port")
	sc := flag.String("sc", "/ip4/127.0.0.1/udp/3000", "STUN server")
	flag.Parse()

	prvKey, _, err := crypto.GenerateKeyPair(crypto.RSA, 2048)
	if err != nil {
		panic(err)
	}
	c, err := net.ListenPacket("udp4", "0.0.0.0:"+*port)
	if err != nil {
		panic(err)
	}

	tlsConfig, err := stun.GenerateConfig(prvKey)
	if err != nil {
		panic(err)
	}
	l, err := quic.Listen(c, tlsConfig, quicConfig)
	if err != nil {
		panic(err)
	}
	go handleListener(l)

	client, err := stun.NewClient(prvKey, c)
	if err != nil {
		panic(err)
	}

	saddr, _ := ma.NewMultiaddr(*sc)
	client.ConnectSTUNServer([]ma.Multiaddr{saddr})

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter dest Multiaddr: ")
	text, _ := reader.ReadString('\n')
	text = strings.Trim(text, "\n")

	fmt.Println("Punching hole to: ", text)

	raddr, err := ma.NewMultiaddr(text)
	if err != nil {
		panic(err)
	}
	com, err := client.PunchHole(raddr)
	if err != nil {
		panic(err)
	}

	<-com

	addr, err := manet.ToNetAddr(raddr)
	time.Sleep(time.Second * 1)

	fmt.Println("Now trying to connect to", addr)
	sess, err := quic.Dial(c, addr, "proto.p2p", tlsConfig, quicConfig)
	if err != nil {
		panic(err)
	}

	stream, err := sess.OpenStream()
	if err != nil {
		panic(err)
	}

	stream.Write([]byte("hello: " + c.LocalAddr().String()))
	for {
		fmt.Print("> ")
		text, _ = reader.ReadString('\n')
		text = strings.Trim(text, "\n")

		stream.Write([]byte(text))
	}
	select {}

}

func handleListener(l quic.Listener) {
	for {
		sess, err := l.Accept()
		if err != nil {
			fmt.Println("Error in accept session:", err)
		}
		stream, err := sess.AcceptStream()
		if err != nil {
			fmt.Println("Error in accept session:", err)
		}
		go handleStream(stream)
	}
}

func handleStream(s quic.Stream) {
	buf := make([]byte, 1000)
	for {
		i, err := s.Read(buf)
		if err != nil {
			fmt.Println("Err when reading", err)
			continue
		}
		fmt.Printf("$ \x1b[32m%s\x1b[0m\n", string(buf[:i]))
	}
}
