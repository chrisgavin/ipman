package main

import (
	"github.com/chrisgavin/ipman/internal/commands"
)

func main() {
	command, err := commands.NewRootCommand()
	if err != nil {
		panic(err)
	}
	command.Run()
}
