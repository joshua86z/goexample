package log

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"time"
)

type Config struct {
	Level    string `yaml:"level"`
	FileName string `yaml:"filename"`
}

// default
var (
	logger, _           = zap.NewProduction()
	Writer    io.Writer = os.Stdout
)

func FileLogger(file string) io.Writer {
	return &lumberjack.Logger{
		Filename:   file,
		MaxSize:    100,
		MaxAge:     28,
		MaxBackups: 10,
		LocalTime:  true,
		Compress:   false,
	}
}

func Init(config Config) {

	var lv zapcore.Level
	err := lv.UnmarshalText([]byte(config.Level))
	if err != nil {
		panic(err)
	}

	ec := zapcore.EncoderConfig{
		TimeKey:       "ts",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "caller",
		FunctionKey:   zapcore.OmitKey,
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.CapitalColorLevelEncoder,
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006/01/02 15:03:04-07"))
		},
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	if config.FileName != "" {
		Writer = FileLogger(config.FileName)

		ec.EncodeLevel = zapcore.CapitalLevelEncoder
	}

	sampler := zap.WrapCore(func(core zapcore.Core) zapcore.Core {

		return zapcore.NewSamplerWithOptions(core, time.Second, 100, 100)
	})

	core := zapcore.NewCore(zapcore.NewConsoleEncoder(ec), zapcore.AddSync(Writer), lv)
	logger = zap.New(core, sampler, zap.AddCaller(), zap.AddCallerSkip(1))

	//logger, _ := zap.NewProduction()

	//config := zap.Config{
	//	Level:       zap.NewAtomicLevelAt(zapcore.InfoLevel),
	//	Development: false,
	//	Sampling: &zap.SamplingConfig{
	//		Initial:    100,
	//		Thereafter: 100,
	//		Hook: func(e zapcore.Entry, s zapcore.SamplingDecision) {
	//		},
	//	},
	//	Encoding: "console",
	//	EncoderConfig: zapcore.EncoderConfig{
	//		TimeKey:        "ts",
	//		LevelKey:       "level",
	//		NameKey:        "logger",
	//		CallerKey:      "caller",
	//		FunctionKey:    zapcore.OmitKey,
	//		MessageKey:     "msg",
	//		StacktraceKey:  "stacktrace",
	//		LineEnding:     zapcore.DefaultLineEnding,
	//		EncodeLevel:    zapcore.LowercaseLevelEncoder,
	//		EncodeTime:     zapcore.EpochTimeEncoder,
	//		EncodeDuration: zapcore.SecondsDurationEncoder,
	//		EncodeCaller:   zapcore.ShortCallerEncoder,
	//	},
	//	OutputPaths:      []string{"stdout"},
	//	ErrorOutputPaths: []string{"stderr"},
	//}
	//
	//logger, _ = config.Build()
}

func Logger() *zap.Logger {
	return logger
}

// Debug logs a message at DebugLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Debug(format string, v ...interface{}) {
	logger.Debug(fmt.Sprintf(format, v...))
}

// Info logs a message at InfoLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Info(format string, v ...interface{}) {
	logger.Info(fmt.Sprintf(format, v...))
}

// Warn logs a message at WarnLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Warn(format string, v ...interface{}) {
	logger.Warn(fmt.Sprintf(format, v...))
}

// Error logs a message at ErrorLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Error(format string, v ...interface{}) {
	logger.Error(fmt.Sprintf(format, v...))
}

// DPanic logs a message at DPanicLevel. The message includes any fields
// passed at the log site, as well as any fields accumulated on the logger.
//
// If the logger is in development mode, it then panics (DPanic means
// "development panic"). This is useful for catching errors that are
// recoverable, but shouldn't ever happen.
func DPanic(format string, v ...interface{}) {
	logger.DPanic(fmt.Sprintf(format, v...))
}

// Panic logs a message at PanicLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
//
// The logger then panics, even if logging at PanicLevel is disabled.
func Panic(format string, v ...interface{}) {
	logger.Panic(fmt.Sprintf(format, v...))
}

// Fatal logs a message at FatalLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
//
// The logger then calls os.Exit(1), even if logging at FatalLevel is
// disabled.
func Fatal(format string, v ...interface{}) {
	logger.Fatal(fmt.Sprintf(format, v...))
}
