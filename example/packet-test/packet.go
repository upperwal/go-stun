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
	ups := strings.Split(text, ",")

	addr1, err := net.ResolveUDPAddr("udp4", ups[0])
	if err != nil {
		panic(err)
	}

	addr2, err := net.ResolveUDPAddr("udp4", ups[1])
	if err != nil {
		panic(err)
	}

	for {
		text, _ := reader.ReadString('\n')
		text = strings.Trim(text, "\n")

		var addr net.Addr
		if text == "1" {
			addr = addr1
		} else {
			addr = addr2
		}
		_, err := c.WriteTo([]byte("hello"), addr)
		if err != nil {
			panic(err)
		}

	}
}
