package main

import (
	"Mxx/api"
	"Mxx/api/log"
	"Mxx/cmd"
	"Mxx/contexts"
	"embed"
	"errors"
	"flag"
	"fmt"
	"go.uber.org/zap"
	"os"
	"runtime"
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
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			fmt.Println("Failed to sync logger")
		}
	}(logger)
	api.StaticFS = webDist

	// init context variables
	ffmpegPath := "ffmpeg"
	if runtime.GOOS == "windows" {
		ffmpegPath = "./ffmpeg.exe"
	}
	contexts.InitContexts(ffmpegPath)

	err = cmd.Run(options)
	if err != nil {
		_, _ = os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}
}
