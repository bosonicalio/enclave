package globallog

import (
	"log/slog"
	"os"
	"sync"
)

var (
	_globalLoggerOnce sync.Once
	_globalLogger     *slog.Logger
)

// Logger returns the global logger instance, initializing it if necessary.
//
// This logger must be used for all logging operations within enclave internal operations. It is not intended
// for use in external applications or libraries.
func Logger() *slog.Logger {
	_globalLoggerOnce.Do(func() {
		// Initialize the global logger with default options
		_globalLogger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,
		}))
	})
	return _globalLogger
}
