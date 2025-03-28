// Package logger provides a logger for the application
// TODO: need to fix this package
package logger

import (
	"io"
	"log/slog"
	"os"
	"sync"

	"github.com/ctfrancia/buho/internal/model"
)

// Logger provides conditional logging based on environment using slog
type Logger struct {
	env      model.Environment
	mu       sync.Mutex
	logger   *slog.Logger
	handlers []slog.Handler
}

// New creates a new Logger with the specified environment
func New(env model.Environment) *Logger {
	logger := &Logger{
		env:      env,
		handlers: make([]slog.Handler, 0),
	}

	// Set default outputs based on environment
	if env == model.Development {
		// For development, use text format with debug level
		textHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
		logger.handlers = append(logger.handlers, textHandler)
	} else {
		// For production, use JSON format with info level
		jsonHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
		logger.handlers = append(logger.handlers, jsonHandler)
	}

	// Create a logger with all handlers
	logger.updateLogger()

	return logger
}

// updateLogger creates or updates the internal slog.Logger
// with the current set of handlers
func (l *Logger) updateLogger() {
	l.mu.Lock()
	defer l.mu.Unlock()

	if len(l.handlers) == 0 {
		// Default to stdout if no handlers
		l.logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
		return
	}

	// Use the first handler directly if there's only one
	if len(l.handlers) == 1 {
		l.logger = slog.New(l.handlers[0])
		return
	}

	// Create a group handler for multiple handlers
	l.logger = slog.New(&multiHandler{handlers: l.handlers})
}

// multiHandler implements slog.Handler to send logs to multiple handlers
type multiHandler struct {
	handlers []slog.Handler
}

// Enabled checks if at least one handler is enabled for the given level
func (h *multiHandler) Enabled(ctx slog.Context, level slog.Level) bool {
	for _, handler := range h.handlers {
		if handler.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

// Handle sends the record to all handlers
func (h *multiHandler) Handle(ctx slog.Context, record slog.Record) error {
	for _, handler := range h.handlers {
		if handler.Enabled(ctx, record.Level) {
			if err := handler.Handle(ctx, record); err != nil {
				return err
			}
		}
	}
	return nil
}

// WithAttrs returns a new handler with the given attributes
func (h *multiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newHandlers := make([]slog.Handler, len(h.handlers))
	for i, handler := range h.handlers {
		newHandlers[i] = handler.WithAttrs(attrs)
	}
	return &multiHandler{handlers: newHandlers}
}

// WithGroup returns a new handler with the given group
func (h *multiHandler) WithGroup(name string) slog.Handler {
	newHandlers := make([]slog.Handler, len(h.handlers))
	for i, handler := range h.handlers {
		newHandlers[i] = handler.WithGroup(name)
	}
	return &multiHandler{handlers: newHandlers}
}

// AddHandler adds a slog.Handler to the logger
func (l *Logger) AddHandler(handler slog.Handler) {
	l.mu.Lock()
	l.handlers = append(l.handlers, handler)
	l.mu.Unlock()

	l.updateLogger()
}

// AddFileHandler adds a file output with JSON formatting
func (l *Logger) AddFileHandler(file *os.File, level slog.Level) error {
	handler := slog.NewJSONHandler(file, &slog.HandlerOptions{
		Level: level,
	})

	l.AddHandler(handler)
	return nil
}

// SetEnvironment changes the logger environment
func (l *Logger) SetEnvironment(env model.Environment) {
	if l.env == env {
		return
	}

	l.mu.Lock()
	l.env = env
	l.mu.Unlock()

	// Adjust log levels based on environment
	if env == model.Development {
		for i, handler := range l.handlers {
			// Try to set debug level for development
			if textHandler, ok := handler.(*slog.TextHandler); ok {
				l.handlers[i] = slog.NewTextHandler(textHandler.Writer(), &slog.HandlerOptions{
					Level: slog.LevelDebug,
				})
			} else if jsonHandler, ok := handler.(*slog.JSONHandler); ok {
				l.handlers[i] = slog.NewJSONHandler(jsonHandler.Writer(), &slog.HandlerOptions{
					Level: slog.LevelDebug,
				})
			}
		}
	} else {
		// For production, set minimum level to info
		for i, handler := range l.handlers {
			if textHandler, ok := handler.(*slog.TextHandler); ok {
				l.handlers[i] = slog.NewTextHandler(textHandler.Writer(), &slog.HandlerOptions{
					Level: slog.LevelInfo,
				})
			} else if jsonHandler, ok := handler.(*slog.JSONHandler); ok {
				l.handlers[i] = slog.NewJSONHandler(jsonHandler.Writer(), &slog.HandlerOptions{
					Level: slog.LevelInfo,
				})
			}
		}
	}

	l.updateLogger()
}

// WithContext returns a logger with additional context fields
func (l *Logger) WithContext(ctx map[string]interface{}) *Logger {
	attrs := make([]any, 0, len(ctx)*2)
	for k, v := range ctx {
		attrs = append(attrs, k, v)
	}

	newLogger := &Logger{
		env:      l.env,
		handlers: l.handlers,
	}
	newLogger.logger = l.logger.With(attrs...)

	return newLogger
}

// Debug logs a debug level message (only shows in Development by default)
func (l *Logger) Debug(msg string, args ...any) {
	l.logger.Debug(msg, args...)
}

// Info logs an info level message
func (l *Logger) Info(msg string, args ...any) {
	l.logger.Info(msg, args...)
}

// Warn logs a warning level message
func (l *Logger) Warn(msg string, args ...any) {
	l.logger.Warn(msg, args...)
}

// Error logs an error level message
func (l *Logger) Error(msg string, args ...any) {
	l.logger.Error(msg, args...)
}

// Fatal logs a fatal error message and exits the program
func (l *Logger) Fatal(msg string, args ...any) {
	l.logger.Error(msg, append(args, "FATAL", true)...)
	os.Exit(1)
}

// GetSlogLogger returns the underlying slog.Logger
// This allows direct access to the slog instance if needed
func (l *Logger) GetSlogLogger() *slog.Logger {
	return l.logger
}

// Writer returns an interface that can be used with slog directly
type Writer interface {
	io.Writer
}
