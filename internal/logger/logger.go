package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	baseLogger   *zap.Logger
	sugarLogger  *zap.SugaredLogger
)

// InitLogger 初始化日志
func InitLogger(debug bool) error {
	// 设置日志级别
	level := zapcore.InfoLevel
	if debug {
		level = zapcore.DebugLevel
	}

	// 编码器配置
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 核心
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.AddSync(os.Stdout),
		level,
	)

	baseLogger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	sugarLogger = baseLogger.Sugar()

	return nil
}

// GetZapLogger 专门给 Gin 框架中间件调用
func GetZapLogger() *zap.Logger {
	return baseLogger
}

// --- 封装常用方法 ---

// Debug 级别
func Debug(args ...interface{}) {
	sugarLogger.Debug(args...)
}

func Debugf(template string, args ...interface{}) {
	sugarLogger.Debugf(template, args...)
}

// Info 级别
func Info(args ...interface{}) {
	sugarLogger.Info(args...)
}

func Infof(template string, args ...interface{}) {
	sugarLogger.Infof(template, args...)
}

// Warn 级别
func Warn(args ...interface{}) {
	sugarLogger.Warn(args...)
}

func Warnf(template string, args ...interface{}) {
	sugarLogger.Warnf(template, args...)
}

// Error 级别
func Error(args ...interface{}) {
	sugarLogger.Error(args...)
}

func Errorf(template string, args ...interface{}) {
	sugarLogger.Errorf(template, args...)
}

// Fatal 级别
func Fatal(args ...interface{}) {
	sugarLogger.Fatal(args...)
}

func Fatalf(template string, args ...interface{}) {
	sugarLogger.Fatalf(template, args...)
}

// With 创建带字段的 logger
func With(fields ...zap.Field) *zap.Logger {
	return baseLogger.With(fields...)
}

// Sync 刷新日志缓冲区
func Sync() error {
	return baseLogger.Sync()
}
