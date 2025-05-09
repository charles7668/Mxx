package converter

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

type M3U8Converter struct {
	ffmpegPath string
}

func CreateM3U8Converter(ffmpegPath string) Converter {
	return &M3U8Converter{ffmpegPath: ffmpegPath}
}

func (converter *M3U8Converter) Convert(input, output string) error {
	if _, err := os.Stat(input); os.IsNotExist(err) {
		return fmt.Errorf("input file %s does not exist", input)
	}
	ffmpegPath := converter.ffmpegPath
	if converter.ffmpegPath == "" {
		if runtime.GOOS == "windows" {
			ffmpegPath = "./ffmpeg"
		} else {
			ffmpegPath = "ffmpeg"
		}
	}

	cmd := exec.Command(
		ffmpegPath,
		"-y",
		"-i", input,
		"-codec", "copy",
		"-start_number", "0",
		"-hls_time", "5",
		"-hls_list_size", "0",
		"-f", "hls",
		output,
	)
	out, err := cmd.CombinedOutput()
	fmt.Println(string(out))
	if err != nil {
		return fmt.Errorf("failed to convert m3u8: %w", err)
	}
	if _, err := os.Stat(output); os.IsNotExist(err) {
		return fmt.Errorf("output file %s does not exist", output)
	}
	return nil
}
