package logging

import (
	"fmt"
	"time"

	"gitlab.com/isteshkov/brute-force-protection/domain/common"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(c *Config) (Logger, error) {
	conf := zap.NewProductionConfig()
	conf.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	conf.EncoderConfig.LevelKey = zapcore.OmitKey
	conf.EncoderConfig.LevelKey = zapcore.OmitKey
	conf.EncoderConfig.TimeKey = zapcore.OmitKey
	conf.EncoderConfig.NameKey = zapcore.OmitKey
	conf.EncoderConfig.CallerKey = zapcore.OmitKey
	conf.EncoderConfig.FunctionKey = zapcore.OmitKey
	conf.EncoderConfig.StacktraceKey = zapcore.OmitKey
	conf.EncoderConfig.LineEnding = zapcore.OmitKey
	conf.EncoderConfig.MessageKey = zapcore.OmitKey
	driver, err := conf.Build()
	if err != nil {
		return nil, err
	}

	fields := map[string]interface{}{
		"instance_id":    c.InstanceId,
		"container_id":   c.ContainerId,
		"container_name": c.ContainerName,
		"env_name":       c.EnvName,
		"version":        c.Version,
		"release":        c.Release,
		"commit_sha":     c.CommitSha,
		"service_name":   c.ServiceName,
	}

	c.logLvl = getLvlFromString(c.LogLvl)
	return &logger{
		fields: fields,
		cfg:    c,
		zap:    driver,
	}, nil
}

type logger struct {
	cfg    *Config
	fields map[string]interface{}
	zap    *zap.Logger
}

func (l logger) WithField(fieldName string, fieldValue interface{}) Logger {
	l.fields = common.CopyStringInterfaceMap(l.fields)
	l.fields[fieldName] = fieldValue
	return &l
}

func (l logger) WithFields(fields map[string]interface{}) Logger {
	l.fields = common.CopyStringInterfaceMap(l.fields)
	for key, value := range fields {
		l.fields[key] = value
	}
	return &l
}

func (l logger) WithContext() Logger {
	panic("implement me")
}

func (l *logger) Debug(msg string, args ...interface{}) {
	if l.cfg.logLvl > lvlDebug {
		return
	}

	fields := common.CopyStringInterfaceMap(l.fields)
	fields[FieldKeyTime] = time.Now().UTC().Format(time.RFC3339Nano)
	fields[FieldKeyLevel] = LevelDebug
	fields[FieldKeyMessage] = fmt.Sprintf(msg, args...)

	l.zap.With(zapcore.Field{
		Key:       "logs",
		Type:      zapcore.ReflectType,
		Interface: fields,
	}).Debug("")
}

func (l *logger) Info(msg string, args ...interface{}) {
	if l.cfg.logLvl > lvlInfo {
		return
	}

	fields := common.CopyStringInterfaceMap(l.fields)
	fields[FieldKeyTime] = time.Now().UTC().Format(time.RFC3339Nano)
	fields[FieldKeyLevel] = LevelInfo
	fields[FieldKeyMessage] = fmt.Sprintf(msg, args...)

	l.zap.With(zapcore.Field{
		Key:       "logs",
		Type:      zapcore.ReflectType,
		Interface: fields,
	}).Debug("")
}

func (l *logger) Warn(msg string, args ...interface{}) {
	if l.cfg.logLvl > lvlWarn {
		return
	}

	fields := common.CopyStringInterfaceMap(l.fields)
	fields[FieldKeyTime] = time.Now().UTC().Format(time.RFC3339Nano)
	fields[FieldKeyLevel] = LevelWarning
	fields[FieldKeyMessage] = fmt.Sprintf(msg, args...)

	l.zap.With(zapcore.Field{
		Key:       "logs",
		Type:      zapcore.ReflectType,
		Interface: fields,
	}).Debug("")
}

type stacktracer interface {
	Stacktrace() string
}

func (l *logger) Error(err error) {
	if l.cfg.logLvl > lvlError {
		return
	}

	fields := common.CopyStringInterfaceMap(l.fields)
	fields[FieldKeyTime] = time.Now().UTC().Format(time.RFC3339Nano)
	fields[FieldKeyLevel] = LevelError
	fields[FieldKeyError] = err.Error()
	fields[FieldKeyMessage] = "Error occurred: " + err.Error()
	if s, ok := err.(stacktracer); ok {
		fields[FieldKeyStacktrace] = s.Stacktrace()
	}

	l.zap.With(zapcore.Field{
		Key:       "logs",
		Type:      zapcore.ReflectType,
		Interface: fields,
	}).Debug("")
}

func (l *logger) ErrorF(err error, msg string, args ...interface{}) {
	if l.cfg.logLvl > lvlError {
		return
	}

	fields := common.CopyStringInterfaceMap(l.fields)
	fields[FieldKeyTime] = time.Now().UTC().Format(time.RFC3339Nano)
	fields[FieldKeyLevel] = LevelError
	fields[FieldKeyError] = err.Error()
	fields[FieldKeyMessage] = "Error occurred: " + fmt.Sprintf(msg, args...)
	if s, ok := err.(stacktracer); ok {
		fields[FieldKeyStacktrace] = s.Stacktrace()
	}

	l.zap.With(zapcore.Field{
		Key:       "logs",
		Type:      zapcore.ReflectType,
		Interface: fields,
	}).Debug("")
}

func (l *logger) Fatal(msg string, args ...interface{}) {
	if l.cfg.logLvl > lvlFatal {
		return
	}

	fields := common.CopyStringInterfaceMap(l.fields)
	fields[FieldKeyTime] = time.Now().UTC().Format(time.RFC3339Nano)
	fields[FieldKeyLevel] = LevelFatal
	fields[FieldKeyMessage] = "Error occurred: " + fmt.Sprintf(msg, args...)

	l.zap.With(zapcore.Field{
		Key:       "logs",
		Type:      zapcore.ReflectType,
		Interface: fields,
	}).Debug("")
}
