package main

import (
	logging "github.com/ipfs/go-log"
	crypto "github.com/libp2p/go-libp2p-crypto"
	stun "github.com/upperwal/go-stun"
)

func main() {

	logging.SetLogLevel("stun", "DEBUG")

	/* port := flag.String("p", "0", "listener port")
	flag.Parse() */

	prvKey, _, err := crypto.GenerateKeyPair(crypto.RSA, 2048)
	if err != nil {
		panic(err)
	}

	_, err = stun.NewServer(prvKey)
	if err != nil {
		panic(err)
	}

	select {}
}
