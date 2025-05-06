package log

import "go.uber.org/zap"

var (
	logger *zap.Logger
)

func init() {
	nLogger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	logger = nLogger
}

func GetLogger() *zap.Logger {
	return logger
}
