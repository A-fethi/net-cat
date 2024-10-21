package netcat

import (
	"log"
	"os"
)

func LoadChatLogo(filename string) (string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	return string(content), nil
}
