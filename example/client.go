package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"

	crypto "github.com/libp2p/go-libp2p-crypto"
	ma "github.com/multiformats/go-multiaddr"
	stun "github.com/upperwal/go-stun"
)

func main() {

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
	client, err := stun.NewClient(prvKey, c)
	if err != nil {
		panic(err)
	}

	saddr, _ := ma.NewMultiaddr(*sc)
	client.ConnectSTUNServer([]ma.Multiaddr{saddr})

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter dest Multiaddr: ")
	text, _ := reader.ReadString('\n')
	text = strings.Trim(text, "\n")

	fmt.Println("Punching hole to: ", text)

	raddr, err := ma.NewMultiaddr(text)
	if err != nil {
		panic(err)
	}
	err = client.PunchHole(raddr)
	if err != nil {
		panic(err)
	}

}
