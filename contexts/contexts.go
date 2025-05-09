package contexts

import "Mxx/ffmpeg"

var (
	FFMpegInstance *ffmpeg.FFMpeg
)

func InitContexts(ffmpegPath string) {
	FFMpegInstance = ffmpeg.GetInstance(ffmpegPath)
}
