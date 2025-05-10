package ffmpeg_test

import (
	"Mxx/contexts"
	"Mxx/ffmpeg/models"
	_ "Mxx/tests_init"
	"os"
	"path/filepath"
	"testing"
)

func TestFFMpeg_GetSilentSegments(t *testing.T) {
	testDir, exist := os.LookupEnv("FFMPEG_TEST_DIR")
	if !exist {
		t.Fatalf("Please set the FFMPEG_TEST_DIR environment variable to the test directory")
	}
	ffmpegInstance := contexts.FFMpegInstance
	testVideo := filepath.Join(testDir, "test_ffmpeg.mp4")
	analyzeOptions := models.GetDefaultSilentAnalyzeOptions(testVideo)
	silentSegments, err := ffmpegInstance.GetSilentSegments(analyzeOptions)
	if err != nil {
		t.Fatalf("GetSilentSegments() returned an error: %v", err)
	}
	if len(silentSegments) == 0 {
		t.Fatalf("GetSilentSegments() returned no silent segments")
	}
	t.Logf("GetSilentSegments() found %d silent segments", len(silentSegments))
}
