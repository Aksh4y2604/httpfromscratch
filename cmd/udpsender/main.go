package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {

	addr, err := net.ResolveUDPAddr("udp", ":42069")
	if err != nil {
		log.Fatal("error", "error", err)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatal("error", "error", err)
	}

	defer conn.Close()

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Println(">")
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal("error", "error", err)
		}
		conn.Write([]byte(line))
	}

}
