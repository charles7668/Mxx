package cmd

import (
	"Mxx/whisper/downloder"
	"Mxx/whisper/transcription"
	"context"
	"errors"
	"fmt"
	"github.com/ggerganov/whisper.cpp/bindings/go/pkg/whisper"
	"os"
	"path/filepath"
	"time"
)

func Run(options RunOptions) error {
	if filepath.Ext(options.inputFile) != ".wav" {
		return errors.New("input file must be a .wav file")
	}
	backgroundCtx := context.Background()
	if _, err := os.Stat("tiny.en.bin"); os.IsNotExist(err) {
		fmt.Println("Model file not found, start downloading...")
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		ctx, cancel := context.WithCancel(backgroundCtx)
		var downloadErr error
		err = downloder.Download(ctx, "tiny.en", cwd, func(progress float32, err error) {
			if err != nil {
				downloadErr = err
				cancel()
				return
			}
			fmt.Printf("Download Progress: %.2f%%\n", progress)
			if progress >= 100 {
				cancel()
			}
		})
		if err != nil {
			return fmt.Errorf("failed to download model %s : %v", "tiny.en", err)
		}
		<-ctx.Done()
		if downloadErr != nil {
			return fmt.Errorf("failed to download model %s : %v", "tiny.en", downloadErr)
		}
	}
	transcriptionOptions := transcription.CreateOptions()
	transcriptionOptions.Language = "auto"
	transcriptionOptions.ProgressCallback = func(progress int) {
		fmt.Printf("Transcription Progress: %d\n", progress)
	}
	buffer := make([]string, 0)
	transcriptionOptions.SegmentCallback = func(segment whisper.Segment) {
		buffer = append(buffer, fmt.Sprintf("[%6s -> %6s] %s", segment.Start.Truncate(time.Millisecond), segment.End.Truncate(time.Millisecond), segment.Text))
	}
	ctx, cancel := context.WithCancel(backgroundCtx)
	fmt.Printf("Transcribing file %s...\n", options.inputFile)
	err := transcription.Transcribe(ctx, options.inputFile, "tiny.en", transcriptionOptions)
	cancel()
	if err != nil {
		return fmt.Errorf("failed to transcribe file %s : %v", options.inputFile, err)
	}
	if options.outputFile == "" {
		for _, segment := range buffer {
			fmt.Println(segment)
		}
		return nil
	}
	outputFile, err := os.Create(options.outputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file %s : %v", options.outputFile, err)
	}
	defer outputFile.Close()
	for _, segment := range buffer {
		_, err := outputFile.WriteString(segment + "\n")
		if err != nil {
			return fmt.Errorf("failed to write to output file %s : %v", options.outputFile, err)
		}
	}
	return nil
}
