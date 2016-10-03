package main

import (
	"os"

	"github.com/supergiant/supergiant/pkg/cli"
)

func main() {
	cli.New(cli.Client, os.Stdin).Run(os.Args)
}
