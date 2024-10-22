package netcat

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"
)

type User struct {
	Name string
	Conn net.Conn
}

var (
	chatLogo string
	users    []User
	backUp   []string
)

func HandleClient(conn net.Conn) {
	defer conn.Close()
	var err error
	chatLogo, err = LoadChatLogo("./ressources/welcome.txt")
	if err != nil {
		log.Fatalf("Error loading chat logo: %v", err)
		return
	}
	fmt.Fprint(conn, string([]byte(chatLogo)))
	name, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Fatal(err)
		return
	}
	name = name[:len(name)-1]
	for _, v := range backUp {
		fmt.Fprint(conn, v)
	}
	time := time.Now()
	users = append(users, User{Name: name, Conn: conn})
	broadcastName(name)
	for {
		fmt.Fprintf(conn, "[%d-%.2d-%.2d %.2d:%.2d:%.2d][%s]:", time.Year(), time.Month(), time.Day(), time.Hour(), time.Minute(), time.Second(), name)
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			log.Println("Connection closed:", err)
			return
		}
		broadcastMessage(message, name)
	}
}

func broadcastName(name string) {
	for _, user := range users {
		if user.Name != name {
			fmt.Fprintf(user.Conn, "\n%s has joined our chat...\n", name)
		}
	}
}

func broadcastMessage(message string, senderName string) {
	backUp = append(backUp, fmt.Sprintf("%s %s", senderName, message))
	for _, user := range users {
		if user.Name != senderName {
			fmt.Fprintf(user.Conn, "%s: %s", senderName, message)
		}
	}
}
