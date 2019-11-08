package healthchecksio

import (
	"fmt"
	"log"
)

// Logger defines the pluggable interface used for logging
type Logger interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

// NoOpLogger is the default logger, logs nothing
type NoOpLogger struct{}

// Debugf logs nothing
func (l *NoOpLogger) Debugf(format string, args ...interface{}) {}

// Infof logs nothing
func (l *NoOpLogger) Infof(format string, args ...interface{}) {}

// Errorf logs nothing
func (l *NoOpLogger) Errorf(format string, args ...interface{}) {}

// StandardLogger is the default logger, logs nothing
type StandardLogger struct{}

// Debugf logs nothing
func (l *StandardLogger) Debugf(format string, args ...interface{}) {
	f := fmt.Sprintf("[DEBUG] %s", format)
	log.Printf(f, args...)
}

// Infof logs nothing
func (l *StandardLogger) Infof(format string, args ...interface{}) {
	f := fmt.Sprintf("[INFO] %s", format)
	log.Printf(f, args...)
}

// Errorf logs nothing
func (l *StandardLogger) Errorf(format string, args ...interface{}) {
	f := fmt.Sprintf("[ERROR] %s", format)
	log.Printf(f, args...)
}
