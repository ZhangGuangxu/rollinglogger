package rollinglogger

import (
	"log"
	"os"
	"path/filepath"

	"github.com/ZhangGuangxu/misc"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	Debug = "debug"
	Info  = "info"
	Warn  = "warn"
	Error = "error"
)

func getLevel(level string) zapcore.Level {
	switch level {
	case Debug:
		return zapcore.DebugLevel
	case Info:
		return zapcore.InfoLevel
	case Warn:
		return zapcore.WarnLevel
	case Error:
		return zapcore.ErrorLevel
	default:
		return zapcore.ErrorLevel
	}
}

type Options struct {
	Level    string
	Filename string
	MaxSize  int
	MaxAge   int
}

// NewLogger create a rolling logger by options
func NewLogger(opts Options) *zap.Logger {
	return newRotatingJSONFileLogger(opts)
}

func newRotatingJSONFileLogger(opts Options) *zap.Logger {
	logDir := filepath.Dir(opts.Filename)
	if !misc.IsFileExist(logDir) {
		if err := os.MkdirAll(logDir, 0700); err != nil {
			log.Panicf("Could not create log directory [%s] for file [%v], error [%v]\n",
				logDir, opts.Filename, err)
		}
	}

	// lumberjack.Logger is already safe for concurrent use, so we don't need to lock it.
	sync := zapcore.AddSync(&lumberjack.Logger{
		Filename: opts.Filename,
		MaxSize:  opts.MaxSize,
		MaxAge:   opts.MaxAge,
	})
	level := getLevel(opts.Level)
	jsonEncoder := newJSONEncoder()
	core := zapcore.NewCore(jsonEncoder, sync, level)
	options := []zap.Option{zap.AddStacktrace(zap.ErrorLevel)}
	return zap.New(core, options...)
}

func newJSONEncoder() zapcore.Encoder {
	return zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "lv",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stack",
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	})
}
