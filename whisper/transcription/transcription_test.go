package transcription

import (
	"Mxx/whisper/downloder"
	"context"
	"github.com/ggerganov/whisper.cpp/bindings/go/pkg/whisper"
	"os"
	"path/filepath"
	"testing"
	"time"
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
	segmentCallback := func(segment whisper.Segment) {
		t.Logf("segment: [%6s -> %6s] %s", segment.Start.Truncate(time.Millisecond), segment.End.Truncate(time.Millisecond), segment.Text)
	}
	transcribeOptions := CreateOptions()
	transcribeOptions.ProgressCallback = transcribeProgressCallback
	transcribeOptions.SegmentCallback = segmentCallback
	transcribeCtx, cancelTranscribe := context.WithCancel(ctx)
	testWavFilePath := filepath.Join(testDir, "test_whisper.wav")
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
