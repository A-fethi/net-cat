package netcat

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

// Checks if the name provided is a valid name
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

// Prints message if the name provided is an invalid name
func invalidName(conn net.Conn) {
	fmt.Fprintln(conn, "Invalid name. Name must:")
	fmt.Fprintln(conn, "- Not be empty")
	fmt.Fprintln(conn, "- Only contain letters, numbers, underscore (_), or hyphen (-)")
	fmt.Fprintln(conn, "Please try again.")
}

// Reads the name, checks the name if valid using isValidName function, checks the name if already exists
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
			invalidName(conn)
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
			log.Printf("Client fails to change their name %s: %s", conn.RemoteAddr().String(), "("+name+")")
			continue
		}
		return name, nil
	}
}
