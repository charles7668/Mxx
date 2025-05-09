package ffmpeg

import "Mxx/ffmpeg/converter"

type FFMpeg struct {
	FFMpegPath string
}

var (
	singletone *FFMpeg
)

type ConverterType int

const (
	M3U8Converter ConverterType = iota
	AudioConverter
)

func (f *FFMpeg) CreateConverter(converterType ConverterType) converter.Converter {
	switch converterType {
	case M3U8Converter:
		return converter.CreateM3U8Converter(f.FFMpegPath)
	case AudioConverter:
		return converter.CreateAudioConverter(f.FFMpegPath)
	}
	return nil
}

func GetInstance(ffmpegPath string) *FFMpeg {
	if singletone != nil {
		return singletone
	}
	singletone = &FFMpeg{
		FFMpegPath: ffmpegPath,
	}
	return singletone
}
