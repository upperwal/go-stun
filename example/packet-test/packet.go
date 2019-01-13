package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	port := flag.String("p", "0", "listener port")

	c, err := net.ListenPacket("udp4", "0.0.0.0:"+*port)
	if err != nil {
		panic(err)
	}
	fmt.Println("My addr: ", c.LocalAddr())

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter dest Multiaddr: ")
	text, _ := reader.ReadString('\n')
	text = strings.Trim(text, "\n")

	addr, err := net.ResolveUDPAddr("udp4", text)
	if err != nil {
		panic(err)
	}

	for {
		_, err := c.WriteTo([]byte("hello"), addr)
		if err != nil {
			panic(err)
		}
		reader.ReadString('\n')
	}
}
