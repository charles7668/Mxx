package ffmpeg

import (
	"Mxx/ffmpeg/models"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func parseSilenceLog(log string) []models.SilentSegment {
	var segments []models.SilentSegment
	var current models.SilentSegment
	var inSilence bool

	reStart := regexp.MustCompile(`^\[silencedetect @ [0-9a-fx]+\] silence_start: ([\d.]+)`)
	reEnd := regexp.MustCompile(`^\[silencedetect @ [0-9a-fx]+\] silence_end: ([\d.]+) \| silence_duration: ([\d.]+)`)

	lines := strings.Split(log, "\n")
	for _, line := range lines {
		if match := reStart.FindStringSubmatch(line); len(match) == 2 {
			// silence_start
			startSec, _ := strconv.ParseFloat(match[1], 64)
			current = models.SilentSegment{Start: time.Duration(startSec * float64(time.Second))}
			inSilence = true
		} else if match := reEnd.FindStringSubmatch(line); len(match) == 3 && inSilence {
			// silence_end + duration
			endSec, _ := strconv.ParseFloat(match[1], 64)
			durSec, _ := strconv.ParseFloat(match[2], 64)
			current.End = time.Duration(endSec * float64(time.Second))
			current.Duration = time.Duration(durSec * float64(time.Second))
			segments = append(segments, current)
			inSilence = false
		}
	}
	return segments
}

func (f *FFMpeg) GetSilentSegments(options models.SilentAnalyzeOptions) ([]models.SilentSegment, error) {
	cmd := exec.Command(
		f.FFMpegPath,
		"-i",
		options.InputFilePath,
		"-af",
		fmt.Sprintf("silencedetect=noise=%.1fdB:d=%.3f", options.NoiseDB, options.Duration.Seconds()),
		"-f", "null", "-")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to convert audio: %w", err)
	}
	silentSegments := parseSilenceLog(string(out))
	return silentSegments, nil
}
