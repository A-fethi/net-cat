package netcat

import (
	"log"
	"net"
)

var chatLogo string

// HandleClient manages a single client connection.
func HandleClient(conn net.Conn) {
	defer conn.Close()
	var err error
	chatLogo, err = LoadChatLogo("./ressources/welcome.txt")
	if err != nil {
		log.Fatalf("Error loading chat logo: %v", err)
	}
	conn.Write([]byte(chatLogo))
}
