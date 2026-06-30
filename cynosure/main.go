package main

import (
	"log"

	"github.com/Cynosure700/cynosure/cynosure/internal/cli"
)

func main() {
	if err := cli.Main(); err != nil {
		log.Fatal(err)
	}
}
