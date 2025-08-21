package log

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/neee333ko/log/klog"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	logger zap.Logger
}

type VLogger struct {
	logger zap.Logger
	level  int
}

var (
	std *Logger = New(nil)
	mu  sync.Mutex
)

func Init(opts *Options) {
	mu.Lock()
	defer mu.Unlock()
	std = New(opts)
}

func New(opts *Options) *Logger {
	if opts == nil {
		opts = InitOptions()
	}

	var zapLevel Level
	if err := zapLevel.UnmarshalText([]byte(opts.Level)); err != nil {
		zapLevel = zapcore.InfoLevel
	}

	encodeLevel := zapcore.CapitalLevelEncoder
	if opts.EnableColor && opts.Format == "console" {
		encodeLevel = zapcore.CapitalColorLevelEncoder
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    encodeLevel,
		EncodeTime:     timeEncoder,
		EncodeDuration: milliSecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	config := zap.Config{
		Level:       zap.NewAtomicLevelAt(zapLevel),
		Development: opts.Development,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:          opts.Format,
		EncoderConfig:     encoderConfig,
		OutputPaths:       opts.OutputPaths,
		ErrorOutputPaths:  opts.ErrorOutputPaths,
		DisableCaller:     opts.DisableCaller,
		DisableStacktrace: opts.DisableStacktrace,
	}

	logger, err := config.Build(zap.AddStacktrace(zapcore.PanicLevel), zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}

	Logger := &Logger{
		logger: *logger,
	}

	klog.InitKlog(logger)
	zap.RedirectStdLog(logger)

	return Logger
}

func VtoZapLevel(level int) zapcore.Level {
	if level > 127 {
		level = 127
	}

	if level < 0 {
		level = 0
	}

	return zapcore.Level(0 - level)
}

func V(level int) *VLogger {
	return &VLogger{
		logger: std.logger,
		level:  level,
	}
}

func (l *Logger) V(level int) *VLogger {
	return &VLogger{
		logger: l.logger,
		level:  level,
	}
}

func WithNamed(name string) *Logger {
	return std.WithNamed(name)
}

func (l *Logger) WithNamed(name string) *Logger {
	logger := l.logger.Named(name)

	return &Logger{
		logger: *logger,
	}
}

func WithValues(keysAndValues ...interface{}) *Logger {
	return std.WithValues(keysAndValues...)
}

func (l *Logger) WithValues(keysAndValues ...interface{}) *Logger {
	logger := l.logger.With(handleFields(&l.logger, keysAndValues)...)

	return &Logger{
		logger: *logger,
	}
}

func Flush() error {
	return std.Flush()
}

func (l *Logger) Flush() error {
	return l.logger.Sync()
}

func StdErrorLog() *log.Logger {
	logger, err := zap.NewStdLogAt(&std.logger, zapcore.ErrorLevel)
	if err == nil {
		return logger
	}

	return nil
}

func StdInfoLog() *log.Logger {
	logger, err := zap.NewStdLogAt(&std.logger, zapcore.InfoLevel)
	if err == nil {
		return logger
	}

	return nil
}

func handleFields(l *zap.Logger, args []interface{}, additional ...zap.Field) []zap.Field {
	// a slightly modified version of zap.SugaredLogger.sweetenFields
	if len(args) == 0 {
		// fast-return if we have no suggared fields.
		return additional
	}

	// unlike Zap, we can be pretty sure users aren't passing structured
	// fields (since logr has no concept of that), so guess that we need a
	// little less space.
	fields := make([]zap.Field, 0, len(args)/2+len(additional))
	for i := 0; i < len(args); {
		// check just in case for strongly-typed Zap fields, which is illegal (since
		// it breaks implementation agnosticism), so we can give a better error message.
		if _, ok := args[i].(zap.Field); ok {
			l.DPanic("strongly-typed Zap Field passed to logr", zap.Any("zap field", args[i]))

			break
		}

		// make sure this isn't a mismatched key
		if i == len(args)-1 {
			l.DPanic("odd number of arguments passed as key-value pairs for logging", zap.Any("ignored key", args[i]))

			break
		}

		// process a key-value pair,
		// ensuring that the key is a string
		key, val := args[i], args[i+1]
		keyStr, isString := key.(string)
		if !isString {
			// if the key isn't a string, DPanic and stop logging
			l.DPanic(
				"non-string key argument passed to logging, ignoring all later arguments",
				zap.Any("invalid key", key),
			)

			break
		}

		fields = append(fields, zap.Any(keyStr, val))
		i += 2
	}

	return append(fields, additional...)
}

func (vlog *VLogger) Info(msg string, fields ...zap.Field) {
	if ce := vlog.logger.Check(VtoZapLevel(vlog.level), msg); ce != nil {
		ce.Write(fields...)
	}
}

func (vlog *VLogger) Infof(format string, v ...interface{}) {
	if ce := vlog.logger.Check(VtoZapLevel(vlog.level), fmt.Sprintf(format, v...)); ce != nil {
		ce.Write()
	}
}

func (vlog *VLogger) Infow(msg string, keysAndValues ...interface{}) {
	if ce := vlog.logger.Check(VtoZapLevel(vlog.level), msg); ce != nil {
		ce.Write(handleFields(&vlog.logger, keysAndValues)...)
	}
}

func Debug(msg string, fields ...zap.Field) {
	std.logger.Debug(msg, fields...)
}

func (l *Logger) Debug(msg string, fields ...zap.Field) {
	l.logger.Debug(msg, fields...)
}

func Debugf(format string, v ...interface{}) {
	std.logger.Sugar().Debugf(format, v...)
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	l.logger.Sugar().Debugf(format, v...)
}

func Debugw(msg string, keysAndValues ...interface{}) {
	std.logger.Sugar().Debugw(msg, keysAndValues...)
}

func (l *Logger) Debugw(msg string, keysAndValues ...interface{}) {
	l.logger.Sugar().Debugw(msg, keysAndValues...)
}

func Info(msg string, fields ...zap.Field) {
	std.logger.Info(msg, fields...)
}

func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.logger.Info(msg, fields...)
}

func Infof(format string, v ...interface{}) {
	std.logger.Sugar().Infof(format, v...)
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.logger.Sugar().Infof(format, v...)
}

func Infow(msg string, keysAndValues ...interface{}) {
	std.logger.Sugar().Infow(msg, keysAndValues...)
}

func (l *Logger) Infow(msg string, keysAndValues ...interface{}) {
	l.logger.Sugar().Infow(msg, keysAndValues...)
}

func Warn(msg string, fields ...zap.Field) {
	std.logger.Warn(msg, fields...)
}

func (l *Logger) Warn(msg string, fields ...zap.Field) {
	l.logger.Warn(msg, fields...)
}

func Warnf(format string, v ...interface{}) {
	std.logger.Sugar().Warnf(format, v...)
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	l.logger.Sugar().Warnf(format, v...)
}

func Warnw(msg string, keysAndValues ...interface{}) {
	std.logger.Sugar().Warnw(msg, keysAndValues...)
}

func (l *Logger) Warnw(msg string, keysAndValues ...interface{}) {
	l.logger.Sugar().Warnw(msg, keysAndValues...)
}

func Error(msg string, fields ...zap.Field) {
	std.logger.Error(msg, fields...)
}

func (l *Logger) Error(msg string, fields ...zap.Field) {
	l.logger.Error(msg, fields...)
}

func Errorf(format string, v ...interface{}) {
	std.logger.Sugar().Errorf(format, v...)
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.logger.Sugar().Errorf(format, v...)
}

func Errorw(msg string, keysAndValues ...interface{}) {
	std.logger.Sugar().Errorw(msg, keysAndValues...)
}

func (l *Logger) Errorw(msg string, keysAndValues ...interface{}) {
	l.logger.Sugar().Errorw(msg, keysAndValues...)
}

func DPanic(msg string, fields ...zap.Field) {
	std.logger.DPanic(msg, fields...)
}

func (l *Logger) DPanic(msg string, fields ...zap.Field) {
	l.logger.DPanic(msg, fields...)
}

func DPanicf(format string, v ...interface{}) {
	std.logger.Sugar().DPanicf(format, v...)
}

func (l *Logger) DPanicf(format string, v ...interface{}) {
	l.logger.Sugar().DPanicf(format, v...)
}

func DPanicw(msg string, keysAndValues ...interface{}) {
	std.logger.Sugar().DPanicw(msg, keysAndValues...)
}

func (l *Logger) DPanicw(msg string, keysAndValues ...interface{}) {
	l.logger.Sugar().DPanicw(msg, keysAndValues...)
}

func Panic(msg string, fields ...zap.Field) {
	std.logger.Panic(msg, fields...)
}

func (l *Logger) Panic(msg string, fields ...zap.Field) {
	l.logger.Panic(msg, fields...)
}

func Panicf(format string, v ...interface{}) {
	std.logger.Sugar().Panicf(format, v...)
}

func (l *Logger) Panicf(format string, v ...interface{}) {
	l.logger.Sugar().Panicf(format, v...)
}

func Panicw(msg string, keysAndValues ...interface{}) {
	std.logger.Sugar().Panicw(msg, keysAndValues...)
}

func (l *Logger) Panicw(msg string, keysAndValues ...interface{}) {
	l.logger.Sugar().Panicw(msg, keysAndValues...)
}

func Fatal(msg string, fields ...zap.Field) {
	std.logger.Fatal(msg, fields...)
}

func (l *Logger) Fatal(msg string, fields ...zap.Field) {
	l.logger.Fatal(msg, fields...)
}

func Fatalf(format string, v ...interface{}) {
	std.logger.Sugar().Fatalf(format, v...)
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.logger.Sugar().Fatalf(format, v...)
}

func Fatalw(msg string, keysAndValues ...interface{}) {
	std.logger.Sugar().Fatalw(msg, keysAndValues...)
}

func (l *Logger) Fatalw(msg string, keysAndValues ...interface{}) {
	l.logger.Sugar().Fatalw(msg, keysAndValues...)
}

func L(ctx context.Context) *Logger {
	return std.L(ctx)
}

func (l *Logger) L(ctx context.Context) *Logger {

	l = l.clone()

	if requestID := ctx.Value(KeyRequestID); requestID != nil {
		l = l.WithValues(KeyRequestID, requestID)
	}

	if username := ctx.Value(KeyUsername); username != nil {
		l = l.WithValues(KeyUsername, username)
	}

	if watcher := ctx.Value(KeyWatcher); watcher != nil {
		l = l.WithValues(KeyWatcher, watcher)
	}

	return l
}

func (l *Logger) clone() *Logger {
	copy := *l

	return &copy
}
