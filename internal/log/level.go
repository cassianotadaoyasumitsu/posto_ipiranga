package log

type level byte

const (
	levelDebug level = 1 << iota
	levelInfo
	levelWarn
	levelError
)

type levelValue struct {
	name string
	level
}

func (v *levelValue) String() string { return v.name }
func (v *levelValue) levelVal()      {}

// Level is the interface that each of the canonical level values implement.
// It contains unexported methods that prevent types from other packages from
// implementing it and guaranteeing that NewFilter can distinguish the levels
// defined in this package from all other values.
type Level interface {
	String() string
	levelVal()
}

var (
	// ErrorLevel represents the Error level, it can also be used as the unique
	// value added to log events in Error level.
	ErrorLevel = &levelValue{level: levelError, name: "error"}

	// WarnLevel represents the Warn level, it can also be used as the unique
	// value added to log events in Warn level.
	WarnLevel = &levelValue{level: levelWarn, name: "warn"}

	// InfoLevel represents the Info level, it can also be used as the unique
	// value added to log events in Info level.
	InfoLevel = &levelValue{level: levelInfo, name: "info"}

	// DebugLevel represents the Debug level, it can also be used as the unique
	// value added to log events in Debug level.
	DebugLevel = &levelValue{level: levelDebug, name: "debug"}
)
