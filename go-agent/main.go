package main

import (
	"log"

	"nano_cc/internal/cli"
)

func main() {
	if err := cli.Main(); err != nil {
		log.Fatal(err)
	}
}
