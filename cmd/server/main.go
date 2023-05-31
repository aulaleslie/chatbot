package main

import (
	"chatbot/app/server"
	"os"
)

func main() {
	if err := server.Run(); err != nil {
		os.Exit(1)
	}
}
