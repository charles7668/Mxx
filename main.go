package main

import (
	"Mxx/cmd"
	"os"
)

func main() {
	args := os.Args
	options, err := cmd.ParseArgs(args)
	if err != nil {
		_, _ = os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}
	err = cmd.Run(options)
	if err != nil {
		_, _ = os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}
}
