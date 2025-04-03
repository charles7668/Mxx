package transcription

import (
	"Mxx/whisper/downloder"
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestTranscription(t *testing.T) {
	testDir, findEnv := os.LookupEnv("WHISPER_TEST_DIR")
	if !findEnv {
		t.Fatalf("Please set the WHISPER_TEST_DIR environment variable to the test directory")
	}
	ctx := context.Background()
	downloadCtx, cancelDownload := context.WithCancel(ctx)
	var downloadErr error
	downloadErr = nil
	err := downloder.Download(downloadCtx, "tiny.en", testDir, func(progress float32, err error) {
		if progress >= 100 {
			cancelDownload()
		}
		if err != nil {
			downloadErr = err
			cancelDownload()
		}
	})
	if err != nil {
		t.Fatalf("download failed: %v", err)
	}
	if downloadErr != nil {
		t.Fatalf("download failed: %v", downloadErr)
	}
	<-downloadCtx.Done()
	transcribeProgressCallback := func(progress int) {
		t.Logf("transcription progress: %d", progress)
	}
	transcribeOptions := CreateOptions()
	transcribeOptions.progressCallback = transcribeProgressCallback
	transcribeCtx, cancelTranscribe := context.WithCancel(ctx)
	testWavFilePath := filepath.Join(testDir, "jfk.wav")
	testModelFilePath := filepath.Join(testDir, "tiny.en.bin")
	go func() {
		err = Transcribe(transcribeCtx, testWavFilePath, testModelFilePath, transcribeOptions)
		defer cancelTranscribe()
		if err != nil {
			t.Errorf("transcription failed: %v", err)
			return
		}
	}()
	<-transcribeCtx.Done()
}
