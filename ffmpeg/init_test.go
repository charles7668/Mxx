package ffmpeg_test

import (
	"Mxx/contexts"
	"runtime"
)

func init() {
	prepareContexts()
}

func prepareContexts() {
	ffmpegPath := "ffmpeg"
	if runtime.GOOS == "windows" {
		ffmpegPath = "./ffmpeg.exe"
	}
	contexts.InitContexts(ffmpegPath)
}
