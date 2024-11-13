package log

import "go.uber.org/zap"

var sugar *zap.SugaredLogger

func newLogger() *zap.SugaredLogger {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	sugar = logger.Sugar()

	return sugar
}

func Logger() *zap.SugaredLogger {
	if sugar == nil {
		sugar = newLogger()
	}

	return sugar
}
