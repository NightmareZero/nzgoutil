package log

import (
	"github.com/NightmareZero/nzgoutil/common"
	"github.com/NightmareZero/nzgoutil/uos"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	defaultLogger *zap.Logger
)

// filed可以用zap.Int()等创建
func Debug(msg string, fields ...zap.Field) {
	defaultLogger.Debug(msg, fields...)
}

// info
// 比上面的方法慢一倍，热点代码不建议用
func Debugf(msg string, value ...any) {
	defaultLogger.Sugar().Debugf(msg, value...)
}

// filed可以用zap.Int()等创建
func Info(msg string, fields ...zap.Field) {
	defaultLogger.Info(msg, fields...)
}

// 比上面的方法慢一倍，热点代码不建议用
func Infof(msg string, value ...any) {
	defaultLogger.Sugar().Infof(msg, value...)
}

// warn
// filed可以用zap.Int()等创建
func Warn(msg string, fields ...zap.Field) {
	defaultLogger.Warn(msg, fields...)
}

// 比上面的方法慢一倍，热点代码不建议用
func Warnf(msg string, fields ...any) {
	defaultLogger.Sugar().Warnf(msg, fields...)
}

// error
// filed可以用zap.Int()等创建
func Error(msg string, fields ...zap.Field) {
	defaultLogger.Error(msg, fields...)
}

// 比上面的方法慢一倍，热点代码不建议用
func Errorf(msg string, fields ...any) {
	defaultLogger.Sugar().Errorf(msg, fields...)
}

// 如果是开发模式，则抛出异常
func DevPanic(msg string, fields ...zap.Field) {
	defaultLogger.DPanic(msg, fields...)
}
func DevPanicf(msg string, fields ...any) {
	defaultLogger.Sugar().DPanicf(msg, fields...)
}

func Panic(msg string, fields ...zap.Field) {
	defaultLogger.Panic(msg, fields...)
}
func Panicf(msg string, fields ...any) {
	defaultLogger.Sugar().Panicf(msg, fields...)
}

func InitWithConfig(config LogConfig) {
	// init logger encoderConfig
	var eConfig zap.Config
	if config.Dev {
		eConfig = zap.NewDevelopmentConfig()
	} else {
		eConfig = zap.NewProductionConfig()
	}
	eConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	var core zapcore.Core

	if !config.MergeErrorLog {
		core = newMultiCoreLogger(config, eConfig)
	} else {
		core = newSingleCoreLogger(config, eConfig)
	}

	defaultLogger = zap.New(core)
}

func newSingleCoreLogger(config LogConfig, eConfig zap.Config) zapcore.Core {
	fixedPath := uos.FixPathEndSlash(common.If(len(config.Path) > 0, config.Path, "./log/"))

	// init logger output file
	logWriter := zapcore.AddSync(getWriter(config.Sync, fixedPath, "log"))

	return zapcore.NewCore(
		zapcore.NewConsoleEncoder(eConfig.EncoderConfig),
		logWriter, config.Level)
}

func newMultiCoreLogger(config LogConfig, eConfig zap.Config) zapcore.Core {
	fixedPath := uos.FixPathEndSlash(common.If(len(config.Path) > 0, config.Path, "./log/"))

	// init logger output file
	logWriter := zapcore.AddSync(getWriter(config.Sync, fixedPath, "log"))
	fixedErrPath := uos.FixPathEndSlash(common.If(len(config.ErrPath) > 0, config.ErrPath, "./log/"))
	// init error logger output file
	errorLogWriter := zapcore.AddSync(getWriter(config.Sync, fixedErrPath, "err"))
	return zapcore.NewTee(
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(eConfig.EncoderConfig),
			logWriter, LogNormalLevel{config.Level, config.MergeErrorLog}),
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(eConfig.EncoderConfig),
			errorLogWriter, zap.ErrorLevel),
	)
}

func InitLog(sync bool) {
	InitWithConfig(LogConfig{
		Sync:  sync,
		Level: zapcore.InfoLevel,
	})
}

type LogNormalLevel struct {
	Level         zapcore.Level
	MergeErrorLog bool
}

func (e LogNormalLevel) Enabled(lvl zapcore.Level) bool {
	return e.Level <= lvl && (!e.MergeErrorLog || lvl < zap.ErrorLevel)
}

type LogConfig struct {
	Sync          bool
	Path          string
	MergeErrorLog bool
	ErrPath       string
	Level         zapcore.Level
	Dev           bool
}
