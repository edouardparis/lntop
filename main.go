// Released under the MIT License
//
// Lightning is a decentralized network using smart contract functionality
// in the Bitcoin protocol to enable instant payments across a network of
// participants. Precisely bidirectional payments channels are set up
// between participants. For more information: http://lightning.network.
//
// lntop is an interactive text-mode channels viewer for Unix systems.
// It supports for the moment the Go implementation lnd only.
package main

import (
	"log"
	"os"

	"github.com/edouardparis/lntop/cli"
)

const Version = "v0.3.1"

func main() {
	err := cli.New(Version).Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
