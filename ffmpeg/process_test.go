package ffmpeg_test

import (
	"Mxx/contexts"
	"Mxx/ffmpeg/models"
	"github.com/google/uuid"
	"os"
	"path/filepath"
	"testing"
)

func TestFFMpeg_SplitBySilentSegmentsToAudio(t *testing.T) {
	testDir, exist := os.LookupEnv("FFMPEG_TEST_DIR")
	if !exist {
		t.Fatalf("Please set the FFMPEG_TEST_DIR environment variable to the test directory")
	}
	tempUUID, err := uuid.NewUUID()
	if err != nil {
		t.Fatalf("Failed to create UUID: %v", err)
	}
	tempPath := filepath.Join("data/temp", tempUUID.String())
	inputFile := filepath.Join(testDir, "test_ffmpeg.mp4")
	silentAnalyzeOptions := models.GetDefaultSilentAnalyzeOptions(inputFile)
	segments, err := contexts.FFMpegInstance.GetSilentSegments(silentAnalyzeOptions)
	if err != nil {
		t.Fatalf("GetSilentSegments() returned an error: %v", err)
	}
	if err := os.MkdirAll(tempPath, os.ModePerm); err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	//defer os.RemoveAll(tempPath)
	cutSegments, err := contexts.FFMpegInstance.SplitBySilentSegmentsToAudio(segments, inputFile, tempPath)
	if err != nil {
		t.Fatalf("SplitBySilentSegmentsToAudio() returned an error: %v", err)
	}
	t.Logf("SplitBySilentSegmentsToAudio() splitCount %d", len(cutSegments))
}
