package utils

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitZap() {
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, os.Stdout, zapcore.InfoLevel)
	logg := zap.New(core, zap.AddCaller())

	zap.ReplaceGlobals(logg)
	defer logg.Sync()
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.LevelKey = "level"
	encoderConfig.TimeKey = "time"
	encoderConfig.MessageKey = "key"
	encoderConfig.EncodeTime = zapcore.TimeEncoder(func(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		loc, err := time.LoadLocation("Asia/Bangkok")
		if err != nil {
			panic(err)
		}
		encoder.AppendString(t.In(loc).Format("02-Jan-2006 15:04:05"))
	})
	encoderConfig.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	return zapcore.NewConsoleEncoder(encoderConfig)
}
