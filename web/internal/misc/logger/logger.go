package logger

import (
	"fmt"
	"os"
	"path/filepath"

	"quiz_platform/internal/misc/apperrors"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	LOGGER_FILE_NAME = "app.log"
)

var (
	GlobalLogger *zap.SugaredLogger

	logFile *os.File
)

func InitLogger(dirPath string) error {
	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		return &apperrors.ErrInternal{
			Message: fmt.Sprintf("cannot create logging dir: %v", err.Error()),
		}
	}

	logFile, err = os.OpenFile(
		filepath.Join(dirPath, LOGGER_FILE_NAME),
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0666)
	if err != nil {
		return &apperrors.ErrInternal{
			Message: fmt.Sprintf("cannot create logging file: %v", err.Error()),
		}
	}

	encoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		MessageKey:    "message",
		LevelKey:      "level",
		TimeKey:       "time",
		CallerKey:     "caller",
		StacktraceKey: "stacktrace",
		EncodeTime:    zapcore.ISO8601TimeEncoder,
		EncodeLevel:   zapcore.LowercaseLevelEncoder,
		EncodeCaller:  zapcore.ShortCallerEncoder,
	})

	core := zapcore.NewCore(encoder, zapcore.AddSync(logFile), zapcore.DebugLevel)

	GlobalLogger = zap.New(core, zap.AddCaller()).Sugar()

	return nil
}

func CleanLogger() {
	defer logFile.Close()
	defer GlobalLogger.Sync()
}
