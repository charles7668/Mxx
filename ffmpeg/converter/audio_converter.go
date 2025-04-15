package converter

import (
	"fmt"
	"os"
	"os/exec"
)

type AudioConverter struct {
	ffmpegPath string
}

func CreateAudioConverter(ffmpegPath string) *AudioConverter {
	return &AudioConverter{ffmpegPath: ffmpegPath}
}

func (converter *AudioConverter) Convert(input, output string) error {
	if _, err := os.Stat(input); os.IsNotExist(err) {
		return fmt.Errorf("input file %s does not exist", input)
	}
	ffmpegPath := converter.ffmpegPath
	if converter.ffmpegPath == "" {
		ffmpegPath = "ffmpeg"
	}

	cmd := exec.Command(fmt.Sprintf("%s -i %s -ar 16000 -ac 1 -c:a pcm_s16le %s", ffmpegPath, input, output))
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to convert audio: %w", err)
	}
	if _, err := os.Stat(output); os.IsNotExist(err) {
		return fmt.Errorf("input file %s does not exist", input)
	}
	return nil
}
