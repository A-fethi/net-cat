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
	users = append(users, User{Name: name, Conn: conn})
	broadcast("name", "", name)
	for {
		// time := time.Now()
		// fmt.Fprintf(conn, "[%d-%.2d-%.2d %.2d:%.2d:%.2d][%s]:", time.Year(), time.Month(), time.Day(), time.Hour(), time.Minute(), time.Second(), name)
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			log.Println("Connection closed:", err)
			return
		}
		broadcast("message", message, name)
	}
}

func broadcast(eventType string, content string, senderName string) {
	time := time.Now()
	for _, user := range users {
		if eventType == "name" && user.Name != senderName {
			fmt.Fprintf(user.Conn, "\n%s has joined our chat...\n", senderName)
		} else if eventType == "message" && user.Name != senderName {
			fmt.Fprintf(user.Conn, "\n[%d-%.2d-%.2d %.2d:%.2d:%.2d][%s]:%s", time.Year(), time.Month(), time.Day(), time.Hour(), time.Minute(), time.Second(), senderName, content)
		}
		fmt.Fprintf(user.Conn, "[%d-%.2d-%.2d %.2d:%.2d:%.2d][%s]:", time.Year(), time.Month(), time.Day(), time.Hour(), time.Minute(), time.Second(), user.Name)
	}

	if eventType == "message" {
		backUp = append(backUp, fmt.Sprintf("[%d-%.2d-%.2d %.2d:%.2d:%.2d][%s]:%s", time.Year(), time.Month(), time.Day(), time.Hour(), time.Minute(), time.Second(), senderName, content))
	}
}
