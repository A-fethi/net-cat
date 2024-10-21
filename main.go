package main

import (
	"fmt"
	"log"
	"net"
	"os"

	netcat "netcat/ressources"
)

const defaultPort = "8989"

func main() {
	port := defaultPort
	if len(os.Args) == 2 {
		port = os.Args[1]
	} else if len(os.Args) > 2 {
		fmt.Println("[USAGE]: ./TCPChat $port")
		return
	}

	listener, err := net.Listen("tcp", ":"+port)
	fmt.Println("Listening on localhost :" + port)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go netcat.HandleClient(conn)
	}
}
