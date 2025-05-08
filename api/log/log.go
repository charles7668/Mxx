package log

import "go.uber.org/zap"

var (
	apiLogger   *zap.Logger
	innerLogger *zap.Logger
)

func init() {
	nLogger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	nInnerLogger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	apiLogger = nLogger
	innerLogger = nInnerLogger
}

func GetApiLogger() *zap.Logger {
	return apiLogger
}

func GetInnerLogger() *zap.Logger {
	return innerLogger
}
