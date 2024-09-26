package main

import (
	"github.com/gabefiori/ts/cli"
	"log"
)

func main() {
	if err := cli.Run(); err != nil {
		log.Fatal(err)
	}
}
