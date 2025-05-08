package cmd

import (
	"Mxx/api"
	"Mxx/api/graceful"
	"Mxx/ffmpeg/converter"
	"Mxx/whisper/downloader"
	"Mxx/whisper/transcription"
	"context"
	"fmt"
	"github.com/ggerganov/whisper.cpp/bindings/go/pkg/whisper"
	"github.com/google/uuid"
	"os"
	"path/filepath"
	"time"
)

func Run(options RunOptions) error {
	if options.webMode {
		router := api.GetApiRouter()
		var routeErr error
		routeCtx, routeCtxCancel := context.WithCancel(graceful.BackgroundContext)
		routeErr = nil
		go func() {
			err := router.Run(":8080")
			if err != nil {
				routeErr = fmt.Errorf("failed to start web server: %v", err)
			}
			routeCtxCancel()
		}()
		<-routeCtx.Done()
		return routeErr
	}
	tempUUID, err := uuid.NewUUID()
	if err != nil {
		return fmt.Errorf("failed to generate uuid for create temp path : %v", err)
	}
	tempPath := tempUUID.String()
	if err := os.MkdirAll(tempPath, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create temp path %s : %v", tempPath, err)
	}
	tempFile := filepath.Join(tempPath, "output.wav")
	// try to convert the media file to wav
	audioConverter := converter.CreateAudioConverter("ffmpeg")
	err = audioConverter.Convert(options.inputFile, tempFile)
	defer os.RemoveAll(tempPath)
	if err != nil {
		return fmt.Errorf("failed to convert audio file %s : %v", options.inputFile, err)
	}
	backgroundCtx := context.Background()
	model := options.model
	modelFile := model
	if filepath.Ext(model) != ".bin" {
		modelFile = model + ".bin"
	}
	if _, err := os.Stat(modelFile); os.IsNotExist(err) {
		fmt.Println("Model file not found, start downloading...")
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		ctx, cancel := context.WithCancel(backgroundCtx)
		var downloadErr error
		err = downloader.Download(ctx, model, cwd, func(progress float32, err error) {
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
			return fmt.Errorf("failed to download model %s : %v", model, err)
		}
		<-ctx.Done()
		if downloadErr != nil {
			return fmt.Errorf("failed to download model %s : %v", model, downloadErr)
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
	err = transcription.Transcribe(ctx, tempFile, model, transcriptionOptions)
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
