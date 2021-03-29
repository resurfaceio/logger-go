package logger

import (
	"strings"
)

type Options struct {
	rules   string
	url     string
	enabled bool
	queue   []string
}

const httpLoggerAgent string = "HttpLogger.go"

//base HttpLogger definition
type HttpLogger struct {
	*BaseLogger
	rules *HttpRules
}

// initialize HttpLogger
func NewHttpLogger(options Options) *HttpLogger {
	baseLogger := NewBaseLogger(httpLoggerAgent, options.url, options.enabled, options.queue)

	loggerRules, _ := newHttpRules(options.rules)

	logger := &HttpLogger{
		baseLogger,
		loggerRules,
	}

	logger.skipCompression = loggerRules.skipCompression
	logger.skipSubmission = loggerRules.skipSubmission

	if (logger.url != "") && (strings.HasPrefix(logger.url, "http:") && !logger.rules.allowHttpUrl) {
		logger.enableable = false
		logger.enabled = false
	}

	return logger
}

// getter for rules
func (logger *HttpLogger) Rules() *HttpRules {
	return logger.rules
}
