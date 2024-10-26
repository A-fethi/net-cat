package netcat

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

const MaxUsers = 10

type User struct {
	Name string
	Conn net.Conn
}

var (
	chatLogo string
	users    []User
	backUp   []string
	mu       sync.RWMutex
)

func isValidName(name string) bool {
	if len(strings.TrimSpace(name)) == 0 {
		return false
	}

	for _, char := range name {
		if !((char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == '_' || char == '-') {
			return false
		}
	}
	return true
}

func readValidName(conn net.Conn) (string, error) {
	reader := bufio.NewReader(conn)
	for {
		fmt.Fprint(conn, "[ENTER YOUR NAME]:")
		name, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		name = strings.TrimSpace(name)
		if !isValidName(name) {
			fmt.Fprintln(conn, "\nInvalid name. Name must:")
			fmt.Fprintln(conn, "- Not be empty")
			fmt.Fprintln(conn, "- Only contain letters, numbers, underscore (_), or hyphen (-)")
			fmt.Fprintln(conn, "Please try again.")
			continue
		}
		mu.Lock()
		nameExists := false
		for _, user := range users {
			if strings.EqualFold(user.Name, name) {
				nameExists = true
				break
			}
		}
		mu.Unlock()
		if nameExists {
			fmt.Fprintln(conn, "This name is already taken. Please choose another name.")
			continue
		}
		return name, nil
	}
}

func removeUser(name string) {
	mu.Lock()
	defer mu.Unlock()
	for i, user := range users {
		if user.Name == name {
			users = append(users[:i], users[i+1:]...)
			break
		}
	}
}

func broadcast(eventType string, content string, senderName string) {
	time := time.Now()
	mu.Lock()
	defer mu.Unlock()
	for _, user := range users {
		if eventType == "name" && user.Name != senderName {
			fmt.Fprintf(user.Conn, "\n%s has joined our chat...\n", senderName)
		} else if eventType == "message" && user.Name != senderName {
			fmt.Fprintf(user.Conn, "\n[%d-%.2d-%.2d %.2d:%.2d:%.2d][%s]:%s", time.Year(), time.Month(), time.Day(), time.Hour(), time.Minute(), time.Second(), senderName, content)
		} else if eventType == "leave" && user.Name != senderName {
			fmt.Fprintf(user.Conn, "\n%s has left the chat\n", senderName)
		}
		fmt.Fprintf(user.Conn, "[%d-%.2d-%.2d %.2d:%.2d:%.2d][%s]:", time.Year(), time.Month(), time.Day(), time.Hour(), time.Minute(), time.Second(), user.Name)
	}

	if eventType == "message" && content != "" {
		backUp = append(backUp, fmt.Sprintf("[%d-%.2d-%.2d %.2d:%.2d:%.2d][%s]:%s", time.Year(), time.Month(), time.Day(), time.Hour(), time.Minute(), time.Second(), senderName, content))
	}
}

func HandleClient(conn net.Conn) {
	mu.RLock()
	if len(users) >= MaxUsers {
		mu.RUnlock()
		fmt.Fprint(conn, "Sorry, the chat room is full (maximum 10 users). Please try again later.\n")
		conn.Close()
		return
	}
	mu.RUnlock()

	defer conn.Close()

	var err error

	chatLogo, err = LoadChatLogo("./ressources/welcome.txt")
	if err != nil {
		log.Fatalf("Error loading chat logo: %v", err)
		return
	}
	fmt.Fprint(conn, string([]byte(chatLogo)))
	name, err := readValidName(conn)
	if err != nil {
		// log.Printf("Error reading name: %v", err)
		return
	}

	mu.RLock()
	for _, v := range backUp {
		fmt.Fprint(conn, v)
	}
	mu.RUnlock()

	mu.Lock()
	users = append(users, User{Name: name, Conn: conn})
	mu.Unlock()

	broadcast("name", "", name)

	defer removeUser(name)
	defer broadcast("leave", "", name)

	for {
		time := time.Now()
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			return
		}
		if message != "\n" {
			// isValid := true

			// for i := 0; i < len(message); i++ {
			// 	if !(message[i] >= 32 && message[i] <= 126) && message[i] != '\n' {
			// 		isValid = false
			// 		continue
			// 	}
			// }
			// if isValid {
			broadcast("message", message, name)
			// } else {
			// 	fmt.Fprint(conn, "Non printable ascii charachters not allowed")
			// 	fmt.Fprintf(conn, "\n[%d-%.2d-%.2d %.2d:%.2d:%.2d][%s]:", time.Year(), time.Month(), time.Day(), time.Hour(), time.Minute(), time.Second(), name)
			// }
		} else if len(message) == 1 {
			fmt.Fprint(conn, "You cannot submit empty message")
			fmt.Fprintf(conn, "\n[%d-%.2d-%.2d %.2d:%.2d:%.2d][%s]:", time.Year(), time.Month(), time.Day(), time.Hour(), time.Minute(), time.Second(), name)
		}
	}
}
