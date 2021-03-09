package main

import (
	"log"
	"os"
)

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		log.Fatal(`Requires a directory argument.`)
	}

	token := gmailToken{}
	token.createToken(args[0])
}
