package converter

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAudioConverter_Convert(t *testing.T) {
	testDir, findEnv := os.LookupEnv("FFMPEG_TEST_DIR")
	if !findEnv {
		t.Fatalf("Please set the FFMPEG_TEST_DIR environment variable to the test directory")
	}
	inputFile := filepath.Join(testDir, "test.mp4")

	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		t.Fatalf("Input file does not exist: %s", inputFile)
	}

	outputFile := filepath.Join(testDir, "output.wav")

	converter := CreateAudioConverter("ffmpeg")

	err := converter.Convert(inputFile, outputFile)
	if err != nil {
		t.Errorf("Convert() returned an error: %v", err)
	}
}

func TestAudioConverter_Convert_InputFileDoesNotExist(t *testing.T) {
	converter := CreateAudioConverter("ffmpeg")

	nonExistentInput := "nonexistent-input.wav"
	outputFile := "output.wav"

	err := converter.Convert(nonExistentInput, outputFile)
	if err == nil {
		t.Errorf("Convert() did not return an error for a non-existent input file")
	}
}
