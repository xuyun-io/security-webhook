package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var Logger *zap.Logger

func init() {
	//logger, err := zap.NewProduction()
	encoderCfg := zapcore.EncoderConfig{
		MessageKey:     "message",
		TimeKey:        "time",
		LevelKey:       "level",
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}
	core := zapcore.NewCore(zapcore.NewJSONEncoder(encoderCfg), zapcore.AddSync(os.Stdout), zap.DebugLevel)

	logger := zap.New(core)

	Logger = logger

}
