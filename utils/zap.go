package utils

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func getFile(f string) *os.File {
	r, err := os.OpenFile(f, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		panic(err)
	}
	return r
}

func InitZap() {
	consoleEncoder, fileEncoder := getEncoder()

	f := getFile("logs/errors.log")

	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, zapcore.AddSync(f), zap.WarnLevel),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zap.InfoLevel),
	)
	logg := zap.New(core, zap.AddCaller())

	zap.ReplaceGlobals(logg)
	defer logg.Sync()
}

func getEncoder() (zapcore.Encoder, zapcore.Encoder) {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.LevelKey = "level"
	encoderConfig.TimeKey = "time"
	encoderConfig.MessageKey = "key"
	encoderConfig.EncodeTime = zapcore.TimeEncoder(func(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		loc, err := time.LoadLocation("Asia/Bangkok")
		if err != nil {
			zap.S().Warn(err.Error())
		}
		encoder.AppendString(t.In(loc).Format("02-Jan-2006 15:04:05"))
	})
	encoderConfig.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
	fileEncoder := zapcore.NewJSONEncoder(encoderConfig)

	return consoleEncoder, fileEncoder
}
