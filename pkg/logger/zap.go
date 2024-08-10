package logger

import (
	"log"
	"os"

	"github.com/streamingfast/logging"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/crypto/ssh/terminal"
)

func newZapLogger(namespace, level, logPath string) *zap.Logger {
	globalLevel := parseLevel(level)

	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})

	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= globalLevel && lvl < zapcore.ErrorLevel
	})

	fileEncoderConfig := zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		TimeKey:        "ts",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
	}
	fileEncoder := zapcore.NewJSONEncoder(fileEncoderConfig)

	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}

	logFileWriter := zapcore.Lock(file)
	logStdErrorWriter := zapcore.Lock(os.Stderr)
	logStdInfoWriter := zapcore.Lock(os.Stdout)

	isTTY := terminal.IsTerminal(int(os.Stderr.Fd()))

	core := zapcore.NewTee(
		zapcore.NewCore(logging.NewEncoder(4, isTTY), logStdErrorWriter, highPriority),
		zapcore.NewCore(logging.NewEncoder(4, isTTY), logStdInfoWriter, lowPriority),
		zapcore.NewCore(fileEncoder, logFileWriter, lowPriority),
	)

	logger := zap.New(
		core,
		zap.AddCaller(), zap.AddCallerSkip(1),
		// zap.AddStacktrace(globalLevel),
	)

	logger = logger.Named(namespace)

	zap.RedirectStdLog(logger)

	return logger
}

func parseLevel(level string) zapcore.Level {
	switch level {
	case LevelDebug:
		return zapcore.DebugLevel
	case LevelInfo:
		return zapcore.InfoLevel
	case LevelWarn:
		return zapcore.WarnLevel
	case LevelError:
		return zapcore.ErrorLevel
	case LevelDPanic:
		return zapcore.DPanicLevel
	case LevelPanic:
		return zapcore.PanicLevel
	case LevelFatal:
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}
