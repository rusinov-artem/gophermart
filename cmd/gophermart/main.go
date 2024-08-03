package main

import (
	"os"

	"github.com/rusinov-artem/gophermart/cmd/gophermart/command"
)

func main() {
	err := command.RootCmd().Execute()
	if err != nil {
		os.Exit(1)
	}
}
