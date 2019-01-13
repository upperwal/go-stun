package main

import (
	"flag"
	"fmt"
	"net"
)

func main() {
	port := flag.String("p", "0", "listener port")

	c, err := net.ListenPacket("udp4", "0.0.0.0:"+*port)
	if err != nil {
		panic(err)
	}
	fmt.Println("My addr: ", c.LocalAddr())

	buf := make([]byte, 1000)

	for {
		i, addr, err := c.ReadFrom(buf)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println("Read: ", buf[:i], "from", addr)
	}
}
