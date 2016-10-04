package main

import (
	"os"

	"github.com/supergiant/supergiant/pkg/cli"
)

var version = "unversioned"

func main() {
	cli.New(cli.Client, os.Stdin, version).Run(os.Args)
}
