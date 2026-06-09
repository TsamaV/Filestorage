package logger

import (
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(logLevel string) (*zap.Logger, func() error, error) {
	lvl := zap.NewAtomicLevel()

	if err := lvl.UnmarshalText([]byte(logLevel)); err != nil {
		return nil, nil, err
	}

	if err := os.MkdirAll("logs", 0755); err != nil {
		return nil, nil, err
	}

	timeStamp := time.Now().UTC().Format("2006-01-02T15-04-05")
	logFilePath := filepath.Join("logs" , timeStamp + ".log")

	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, nil, err
	}

	cfg := zap.NewDevelopmentEncoderConfig()
	cfg.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02T15-04-05")

	encoder := zapcore.NewConsoleEncoder(cfg)

	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), lvl),
		zapcore.NewCore(encoder, zapcore.AddSync(logFile), lvl),
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))

	return logger, logFile.Close, nil
}