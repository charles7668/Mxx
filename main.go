package main

import (
	"Mxx/api"
	"Mxx/api/log"
	"Mxx/cmd"
	"embed"
	"errors"
	"flag"
	"os"
)

//go:embed web/dist/*
var webDist embed.FS

func main() {
	args := os.Args
	options, err := cmd.ParseArgs(args)
	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			os.Exit(0)
		}
		_, _ = os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}
	logger := log.GetApiLogger()
	defer logger.Sync()
	api.StaticFS = webDist
	err = cmd.Run(options)
	if err != nil {
		_, _ = os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}
}
