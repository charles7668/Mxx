package ffmpeg

import (
	"Mxx/ffmpeg/models"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func silentSegmentsToCutSegments(segments []models.SilentSegment) []models.Segment {
	result := make([]models.Segment, 0)
	prevTime := time.Duration(0)
	for i, segment := range segments[0 : len(segments)-1] {
		if segments[i+1].Start-prevTime > 30*time.Second && prevTime < segment.Start {
			newSegment := models.Segment{
				Start:    prevTime,
				End:      segment.Start,
				Duration: segment.Start - prevTime,
			}
			result = append(result, newSegment)
			prevTime = segment.Start
		}
	}
	// insert last segment -1 mean end of file
	{
		newSegment := models.Segment{
			Start:    prevTime,
			End:      -1,
			Duration: -1,
		}
		result = append(result, newSegment)
	}
	return result
}

// SplitBySilentSegmentsToAudio splits the media into segments based on silence detection. output file is named silent_xxx.mp4
func (f *FFMpeg) SplitBySilentSegmentsToAudio(segments []models.SilentSegment, inputMediaFile, outputPath string) ([]models.Segment, error) {
	if len(segments) == 0 {
		return nil, fmt.Errorf("no silent segments found")
	}
	if stat, err := os.Stat(inputMediaFile); os.IsNotExist(err) || stat.IsDir() {
		return nil, fmt.Errorf("input file %s does not exist or is not a file", inputMediaFile)
	}
	if stat, err := os.Stat(outputPath); os.IsNotExist(err) || !stat.IsDir() {
		return nil, fmt.Errorf("output path %s does not exist or is not a directory", outputPath)
	}
	count := 0
	cutSegments := silentSegmentsToCutSegments(segments)
	if len(cutSegments) == 0 {
		return nil, fmt.Errorf("no cut segments found")
	}
	getOutputPath := func(number int) string {
		return filepath.Join(outputPath, fmt.Sprintf("silent_%d.wav", number))
	}
	for _, segment := range cutSegments {
		count++
		outputFile := getOutputPath(count)
		var args []string
		args = append(args,
			"-y",
			"-ss", durationFormat(segment.Start),
			"-i", inputMediaFile,
			"-vn", "-ar", "16000", "-ac", "1", "-c:a", "pcm_s16le")
		if segment.Duration > 0 {
			args = append(args, "-t", durationFormat(segment.Duration))
		}
		args = append(args, outputFile)
		cmd := exec.Command(
			f.FFMpegPath,
			args...,
		)
		_, err := cmd.CombinedOutput()
		if err != nil {
			return nil, fmt.Errorf("failed to split audio: %w", err)
		}
	}

	return cutSegments, nil
}

func durationFormat(duration time.Duration) string {
	// Convert the duration to a string in the format "HH:MM:SS.ms"
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60
	milliseconds := int(duration.Milliseconds()) % 1000
	return fmt.Sprintf("%02d:%02d:%02d.%03d", hours, minutes, seconds, milliseconds)
}
