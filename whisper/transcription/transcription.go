package transcription

import (
	"context"
	"errors"
	"fmt"
	"github.com/ggerganov/whisper.cpp/bindings/go/pkg/whisper"
	"github.com/go-audio/wav"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

const (
	modelExt = ".bin"
)

func getModelNameWithExtension(modelName string) string {
	if filepath.Ext(modelName) == modelExt {
		return modelName
	}
	return modelName + modelExt
}

func Transcribe(ctx context.Context, filePath, modelName string, transcribeOptions TranscribeOptions) error {
	modelFile := getModelNameWithExtension(modelName)
	model, err := whisper.New(modelFile)
	if err != nil {
		return errors.New("failed to create model: " + err.Error())
	}
	file, err := os.Open(filePath)
	if err != nil {
		return errors.New("failed to open file: " + err.Error())
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	var data []float32
	dec := wav.NewDecoder(file)
	if buf, err := dec.FullPCMBuffer(); err != nil {
		return err
	} else if dec.SampleRate != whisper.SampleRate {
		return fmt.Errorf("unsupported sample rate: %d", dec.SampleRate)
	} else if dec.NumChans != 1 {
		return fmt.Errorf("unsupported number of channels: %d", dec.NumChans)
	} else {
		data = buf.AsFloat32Buffer().Data
	}
	whisperContext, err := getWhisperContext(model, transcribeOptions)
	if err != nil {
		return errors.New("failed to get whisper context: " + err.Error())
	}
	encoderCallback := func() bool {
		select {
		case <-ctx.Done():
			log.Println("encoder callback: context done")
			return false
		default:
			return true
		}
	}
	log.Println("starting transcription for file: ", filePath)
	log.Println("using model: ", modelName)
	whisperContext.ResetTimings()
	if err = whisperContext.Process(data, encoderCallback, nil, transcribeOptions.progressCallback); err != nil {
		return errors.New("failed to process audio: " + err.Error())
	}
	whisperContext.PrintTimings()
	if err = Output(whisperContext); err != nil {
		return errors.New("failed to output: " + err.Error())
	}

	return nil
}

func Output(context whisper.Context) error {
	for {
		segment, err := context.NextSegment()
		if err == io.EOF {
			return nil
		} else if err != nil {
			return err
		}
		log.Printf("[%6s->%6s]", segment.Start.Truncate(time.Millisecond), segment.End.Truncate(time.Millisecond))
		log.Println()
		log.Println(segment.Text)
	}
}

func getWhisperContext(model whisper.Model, options TranscribeOptions) (whisper.Context, error) {
	whisperContext, err := model.NewContext()
	if err != nil {
		return nil, errors.New("failed to create whisper context: " + err.Error())
	}
	if model.IsMultilingual() && whisperContext.SetLanguage(options.Language) != nil {
		return nil, errors.New("failed to set language: " + options.Language)
	}
	return whisperContext, nil
}
