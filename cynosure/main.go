package main

import (
	"log"

	"cynosure/internal/cli"
)

func main() {
	if err := cli.Main(); err != nil {
		log.Fatal(err)
	}
}
