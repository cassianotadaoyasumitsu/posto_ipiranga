package zerolog

import (
	"os"

	"git.wealth-park.com/cassiano/posto_ipiranga/internal/log"
	"github.com/rs/zerolog"
)

type Logger struct {
	zlogger zerolog.Logger
}

func toZerologLevel(lvl log.Level) zerolog.Level {
	switch lvl {
	case log.DebugLevel:
		return zerolog.DebugLevel
	case log.InfoLevel:
		return zerolog.InfoLevel
	case log.WarnLevel:
		return zerolog.WarnLevel
	case log.ErrorLevel:
		return zerolog.ErrorLevel
	}
	return zerolog.InfoLevel
}

// mapZerologFields maps the fields in the log package to the zerolog fields
func mapZerologFields() {
	zerolog.TimestampFieldName = log.TimestampFieldName
	zerolog.LevelFieldName = log.LevelFieldName
	zerolog.MessageFieldName = log.MessageFieldName
	zerolog.CallerFieldName = log.CallerFieldName
	zerolog.ErrorFieldName = log.ErrorFieldName
	zerolog.ErrorStackFieldName = log.ErrorStackFieldName
}

// NewConsoleLogger returns an implementation of log.Logger that sends log events to a zerolog.Logger
// and outputs the log events with a zerolog.ConsoleWriter.
func NewConsoleLogger(lvl log.Level) log.Logger {
	mapZerologFields()
	out := zerolog.NewConsoleWriter(
		func(w *zerolog.ConsoleWriter) {
			w.TimeFormat = log.DefaultTimestampLayout
			w.Out = os.Stdout
			w.NoColor = false
		},
	)
	l := &Logger{zerolog.New(out).Level(toZerologLevel(lvl))}
	return log.WithPrefix(l, log.TimestampFieldName, log.DefaultTimestampUTC, log.CallerFieldName, log.Caller(4))
}

// NewJSONLogger returns an implementation of log.Logger that sends log events to a zerolog.Logger
// and outputs the log events in JSON format.
func NewJSONLogger(lvl log.Level) log.Logger {
	mapZerologFields()
	l := &Logger{zerolog.New(os.Stdout).Level(toZerologLevel(lvl))}
	return log.WithPrefix(l, log.TimestampFieldName, log.DefaultTimestampUTC, log.CallerFieldName, log.Caller(4))
}

func (l *Logger) toFields(keyvals ...interface{}) []interface{} {
	if len(keyvals)%2 == 0 {
		return keyvals
	}
	var kvs []interface{}
	for i := 0; i < len(keyvals); i += 2 {
		if i+1 < len(keyvals) {
			kvs = append(kvs, keyvals[i], keyvals[i+1])
		} else {
			kvs = append(kvs, keyvals[i], log.ErrMissingValue)
		}
	}
	return kvs
}

func (l *Logger) Debug(msg string, keyvals ...interface{}) error {
	l.zlogger.Debug().
		Fields(l.toFields(keyvals...)).
		Msg(msg)
	return nil
}

func (l *Logger) Info(msg string, keyvals ...interface{}) error {
	l.zlogger.Info().
		Fields(l.toFields(keyvals...)).
		Msg(msg)
	return nil
}

func (l *Logger) Warn(msg string, keyvals ...interface{}) error {
	l.zlogger.Warn().
		Fields(l.toFields(keyvals...)).
		Msg(msg)
	return nil
}

func (l *Logger) Error(err error, keyvals ...interface{}) error {
	l.zlogger.Error().
		Fields(l.toFields(keyvals...)).
		Err(err).
		Send()
	return nil
}
