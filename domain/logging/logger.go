package logging

type Logger interface {
	WithField(key string, value interface{}) Logger
	WithFields(fields map[string]interface{}) Logger
	WithContext() Logger
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(err error)
	ErrorF(err error, msg string, args ...interface{})
	Fatal(msg string, args ...interface{})
}
