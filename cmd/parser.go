package cmd

import (
	"errors"
	"flag"
)

func ParseArgs(args []string) (RunOptions, error) {
	fs := flag.NewFlagSet("cmd", flag.PanicOnError)
	inputFile := fs.String("i", "", "Path to the input file")
	fs.StringVar(inputFile, "input", "", "Path to the input file (alias for -i)")
	outputFile := fs.String("o", "", "Path to the output file")
	fs.StringVar(outputFile, "output", "", "Path to the output file (alias for -o)")

	err := fs.Parse(args[1:])
	if err != nil {
		return RunOptions{}, err
	}

	if *inputFile == "" {
		return RunOptions{}, errors.New("missing required flag: -input")
	}

	return RunOptions{
		inputFile:  *inputFile,
		outputFile: *outputFile,
	}, nil
}
