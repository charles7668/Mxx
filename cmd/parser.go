package cmd

import (
	"errors"
	"flag"
)

func ParseArgs(args []string) (RunOptions, error) {
	fs := flag.NewFlagSet("cmd", flag.ContinueOnError)
	inputFile := fs.String("i", "", "Path to the input file")
	fs.StringVar(inputFile, "input", "", "Path to the input file (alias for -i)")
	outputFile := fs.String("o", "", "Path to the output file")
	fs.StringVar(outputFile, "output", "", "Path to the output file (alias for -o)")
	model := fs.String("m", "tiny.en", "Whisper model name , you can see all support list in https://github.com/ggml-org/whisper.cpp/blob/master/models/README.md")
	fs.StringVar(model, "model", "tiny.en", "Whisper model name (alias for -m)")
	apiMode := fs.Bool("api", false, "launch api server")
	webMode := fs.Bool("web", false, "launch web server")

	err := fs.Parse(args[1:])
	if err != nil {
		return RunOptions{}, err
	}

	if *inputFile == "" && !*apiMode && !*webMode {
		return RunOptions{}, errors.New("missing required flag: --input/-i or --api or --web")
	}

	return RunOptions{
		inputFile:  *inputFile,
		outputFile: *outputFile,
		model:      *model,
		apiMode:    *apiMode,
		webMode:    *webMode,
	}, nil
}
