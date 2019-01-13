package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strings"
	"time"
)

var (
	list map[string]net.Addr
	c    net.PacketConn
)

func main() {
	port := flag.String("p", "0", "listener port")
	flag.Parse()

	c, err := net.ListenPacket("udp4", "0.0.0.0:"+*port)
	if err != nil {
		panic(err)
	}
	fmt.Println("My addr: ", c.LocalAddr())

	listenCommands(c)

	/* reader := bufio.NewReader(os.Stdin)
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

	} */
}

func listenCommands(c net.PacketConn) {
	reader := bufio.NewReader(os.Stdin)
	list = make(map[string]net.Addr)

	for {
		fmt.Print("$> ")
		text, _ := reader.ReadString('\n')
		text = strings.Trim(text, "\n")
		commandString := strings.SplitN(text, " ", 2)

		switch commandString[0] {
		case "a":
			idx := strings.Split(commandString[1], "-")
			addr, err := net.ResolveUDPAddr("udp4", idx[1])
			if err != nil {
				fmt.Println("Incorrect value: ", err)
				continue
			}
			list[idx[0]] = addr
			fmt.Println("Value Added:", addr)
		case "d":
			dial(c, commandString[1])
		}
	}
}

func dial(c net.PacketConn, s string) {
	addr := list[s]
	if addr == nil {
		fmt.Println("This addr is not added")
		return
	}
	for j := 0; j < 6; j++ {
		for i := 0; i < 3; i++ {
			fmt.Println(i, "Sending packet to:", addr)
			_, err := c.WriteTo([]byte("hello"), addr)
			if err != nil {
				panic(err)
			}
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(300)+20))
		}
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(2000)+1000))
	}

}
