package utils

import (
	"github.com/Calmantara/lis-backend/internal/core/configurations"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	DEFAULT_NAME = "lis-backend"
)

var Log *logger

type (
	Logger interface {
		Debugw(msg string, inputFields ...map[string]interface{})
		Debug(inputFields ...map[string]interface{})
		Errorw(msg string, inputFields ...map[string]interface{})
		Error(inputFields ...map[string]interface{})
		Fatalw(msg string, inputFields ...map[string]interface{})
		Fatal(inputFields ...map[string]interface{})
		Infow(msg string, inputFields ...map[string]interface{})
		Info(inputFields ...map[string]interface{})
	}
	WrappedLoggerInterface interface {
		Debugw(msg string, ops ...interface{})
		Debug(ops ...interface{})
		Errorw(msg string, ops ...interface{})
		Error(ops ...interface{})
		Fatalw(msg string, ops ...interface{})
		Fatal(ops ...interface{})
		Infow(msg string, ops ...interface{})
		Info(ops ...interface{})
	}

	Option func(*logger)

	logger struct {
		wrappedLogger   WrappedLoggerInterface
		environment     configurations.Environment
		cfg             *zap.Config
		opt             []zap.Option
		applicationName string
		initialField    map[string]any
	}
)

func WithLoggerConfig(cfg zap.Config) Option {
	return func(l *logger) {
		l.cfg = &cfg
	}
}

func WithZapOption(opt zap.Option) Option {
	return func(l *logger) {
		l.opt = append(l.opt, opt)
	}
}

func WithAppName(applicationName string) Option {
	return func(l *logger) {
		l.applicationName = applicationName
	}
}

func WithEnvironment(environment configurations.Environment) Option {
	return func(l *logger) {
		l.environment = environment
	}
}

func WithInitialFields(initialFields map[string]any) Option {
	return func(l *logger) {
		l.initialField = initialFields
	}
}

func NewZap(options ...Option) Logger {
	// set default value
	log := &logger{
		environment:     configurations.TEST,
		applicationName: DEFAULT_NAME,
		opt:             []zap.Option{zap.AddCallerSkip(1)},
	}

	// scan all options
	for _, fn := range options {
		fn(log)
	}

	if log.cfg == nil {
		log.cfg = log.newLoggerConfig()
	}
	if len(log.initialField) == 0 {
		log.initialField = log.initialFields(log.applicationName)
	}

	zapLogger, err := log.cfg.Build()
	if err != nil {
		panic(err)
	}

	log.wrappedLogger = zapLogger.WithOptions(log.opt...).Sugar()
	Log = log

	return log
}

func (l *logger) Debugw(msg string, inputFields ...map[string]interface{}) {
	if len(inputFields) > 0 {
		fields := l.transformLogMapToSlice(inputFields[0])
		l.wrappedLogger.Debugw(msg, fields...)

		return
	}

	l.wrappedLogger.Debugw(msg)
}

func (l *logger) Debug(inputFields ...map[string]interface{}) {
	if len(inputFields) > 0 {
		fields := l.transformLogMapToSlice(inputFields[0])
		l.wrappedLogger.Debugw("", fields...)

		return
	}
}

func (l *logger) Errorw(msg string, inputFields ...map[string]interface{}) {
	if len(inputFields) > 0 {
		fields := l.transformLogMapToSlice(inputFields[0])
		l.wrappedLogger.Errorw(msg, fields...)

		return
	}

	l.wrappedLogger.Errorw(msg)
}

func (l *logger) Error(inputFields ...map[string]interface{}) {
	if len(inputFields) > 0 {
		fields := l.transformLogMapToSlice(inputFields[0])
		l.wrappedLogger.Errorw("", fields...)

		return
	}
}

func (l *logger) Fatalw(msg string, inputFields ...map[string]interface{}) {
	if len(inputFields) > 0 {
		fields := l.transformLogMapToSlice(inputFields[0])
		l.wrappedLogger.Fatalw(msg, fields...)

		return
	}

	l.wrappedLogger.Fatalw(msg)
}

func (l *logger) Fatal(inputFields ...map[string]interface{}) {
	if len(inputFields) > 0 {
		fields := l.transformLogMapToSlice(inputFields[0])
		l.wrappedLogger.Fatalw("", fields...)

		return
	}
}

func (l *logger) Infow(msg string, inputFields ...map[string]interface{}) {
	if len(inputFields) > 0 {
		fields := l.transformLogMapToSlice(inputFields[0])
		l.wrappedLogger.Infow(msg, fields...)

		return
	}

	l.wrappedLogger.Infow(msg)
}

func (l *logger) Info(inputFields ...map[string]interface{}) {
	if len(inputFields) > 0 {
		fields := l.transformLogMapToSlice(inputFields[0])
		l.wrappedLogger.Infow("", fields...)

		return
	}
}

func (l *logger) transformLogMapToSlice(logMap map[string]interface{}) []interface{} {
	if len(logMap) == 0 {
		return []interface{}{}
	}

	result := make([]interface{}, len(logMap)*2)
	idx := 0
	for k, v := range logMap {
		result[idx] = k
		result[idx+1] = v
		idx += 2
	}

	return result
}

func (l *logger) newLoggerConfig() *zap.Config {
	encoderCfg := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.EpochTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	cfg := zap.Config{
		Level:             zap.NewAtomicLevelAt(zap.InfoLevel),
		Development:       false,
		DisableCaller:     false,
		DisableStacktrace: false,
		Sampling:          nil,
		Encoding:          "json",
		EncoderConfig:     encoderCfg,
		OutputPaths:       []string{"stdout"},
		ErrorOutputPaths:  []string{"stderr"},
		InitialFields:     l.initialField,
	}

	return &cfg
}

func (l *logger) initialFields(appName string) map[string]interface{} {
	return map[string]interface{}{
		"app":                appName,
		"_filebeat_version":  "2.0.0",
		"_filebeat_selector": "service",
	}
}
