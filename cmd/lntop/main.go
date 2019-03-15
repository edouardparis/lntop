package main

import (
	"log"
	"os"

	"github.com/edouardparis/lntop/cli"
)

func main() {
	err := cli.New().Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
