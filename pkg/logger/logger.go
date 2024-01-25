package logger

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/samber/lo"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type lumberjackSink struct {
	*lumberjack.Logger
}

func (lumberjackSink) Sync() error {
	return nil
}

func defaultCfg() zap.Config {
	encoder := zap.NewProductionEncoderConfig()
	encoder.TimeKey = "time"
	encoder.EncodeTime = zapcore.RFC3339TimeEncoder
	return zap.Config{
		Level:            zap.NewAtomicLevelAt(zap.InfoLevel),
		Development:      false,
		Encoding:         "json",
		EncoderConfig:    encoder,
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}
}

var logger = func() *zap.Logger {
	return lo.Must(defaultCfg().Build(zap.AddCallerSkip(1)))
}()

func ZapLogger() *zap.Logger {
	return logger
}

func SetLogger(level string, dir string, logFile string, dev bool) error {
	l, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}

	rl := &lumberjack.Logger{
		Filename: filepath.Join(dir, logFile),
		MaxSize:  500,
		MaxAge:   7,
		Compress: true,
	}
	if _, err := os.Stat(filepath.Join(dir, logFile)); err == nil {
		// trigger a rotate when service restarts.
		rl.Rotate()
	}
	zap.RegisterSink("lbj", func(*url.URL) (zap.Sink, error) {
		return lumberjackSink{
			Logger: rl,
		}, nil
	})

	cfg := defaultCfg()
	cfg.Level = l
	cfg.OutputPaths = []string{"stderr"}
	if !dev && dir != "" {
		cfg.OutputPaths = append(cfg.OutputPaths, fmt.Sprintf("lbj:%s", logFile))
	}
	logger = lo.Must(cfg.Build(zap.AddCallerSkip(1)))
	return nil
}

func Info(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	logger.Error(msg, fields...)
}

func Debug(msg string, fields ...zap.Field) {
	logger.Debug(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	logger.Fatal(msg, fields...)
}

var (
	Any      = zap.Any
	Duration = zap.Duration
	String   = zap.String
	Strings  = zap.Strings
	Stringer = zap.Stringer
	Int      = zap.Int
	Uint32   = zap.Uint32
	Uint64   = zap.Uint64
	Int64    = zap.Int64
	Float64  = zap.Float64
	Err      = zap.Error
)
