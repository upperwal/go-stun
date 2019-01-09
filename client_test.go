package stun

import (
	"net"
	"testing"
	"time"

	logging "github.com/ipfs/go-log"
	crypto "github.com/libp2p/go-libp2p-crypto"
	ma "github.com/multiformats/go-multiaddr"
)

func createClient(t *testing.T, port string) *Client {
	prvKey, _, err := crypto.GenerateKeyPair(crypto.RSA, 2048)
	if err != nil {
		t.Error(err)
	}
	c, err := net.ListenPacket("udp4", "0.0.0.0:"+port)
	if err != nil {
		t.Error(err)
	}
	client, err := NewClient(prvKey, c)
	if err != nil {
		t.Error(err)
	}

	saddr, _ := ma.NewMultiaddr("/ip4/127.0.0.1/udp/3000")
	client.ConnectSTUNServer([]ma.Multiaddr{saddr})

	return client
}

func TestClient(t *testing.T) {

	logging.SetLogLevel("stun", "DEBUG")

	//ma1, _ := ma.NewMultiaddr("/ip4/127.0.0.1/udp/3001")
	ma2, err := ma.NewMultiaddr("/ip4/127.0.0.1/udp/3002")
	if err != nil {
		t.Fatal(err)
	}

	client1 := createClient(t, "3001")
	createClient(t, "3002")

	time.Sleep(time.Second * 3)

	client1.PunchHole(ma2)

	select {}

}
