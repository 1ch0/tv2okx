package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

var Logger *zap.SugaredLogger

const (
	DebugLevel zapcore.Level = iota - 1
	InfoLevel
)

func SetLevel(level zapcore.Level) {
	customTimeFormat := "2006-01-02 15:04:05.000"

	config := zap.NewProductionEncoderConfig()

	config.EncodeTime = func(time time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(time.Format(customTimeFormat))
	}

	var l *zap.Logger

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(config),
		zapcore.Lock(os.Stdout),
		level,
	)
	switch level {
	case zapcore.DebugLevel:
		l = zap.New(core, zap.AddStacktrace(zap.ErrorLevel), zap.AddCaller())
	default:
		l = zap.New(core, zap.AddCaller())
	}

	Logger = l.Sugar()
}
