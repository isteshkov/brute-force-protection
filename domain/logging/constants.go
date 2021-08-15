package logging

const (
	FieldKeyLevel      = "level"
	FieldKeyError      = "error"
	FieldKeyTime       = "time_utc"
	FieldKeyMessage    = "message"
	FieldKeyRequestID  = "request_id"
	FieldKeyStacktrace = "stacktrace"

	LevelDebug   = "DEBUG"
	LevelInfo    = "INFO"
	LevelWarning = "WARNING"
	LevelError   = "ERROR"
	LevelFatal   = "FATAL"
)
