package logger

import (
	"strings"
)

type Options struct {
	rules   string
	agent   string
	url     string
	enabled bool
	queue   []string
}

const httpLoggerAgent string = "HttpLogger.go"

//base HttpLogger definition
type HttpLogger struct {
	*BaseLogger
	rules HttpRules
}

// initialize HttpLogger
func NewHttpLogger(options Options) *HttpLogger {
	baseLogger := NewBaseLogger(options.agent, options.url, options.enabled, options.queue)

	loggerRules := newHttpRules(options.rules)

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

func (logger *HttpLogger) submitIfPassing(details [][]string) {
	details = logger.rules.apply(details)

	if details == nil {
		return
	}

	details = append(details, []string{"host", logger.host})

	logger.submit(msgStringify(details))
}

// method for converting message details to string format
func msgStringify(msg [][]string) string {
	stringified := ""
	n := len(msg)
	for i, val := range msg {
		stringified += "[" + strings.Join(val, ", ") + "]"
		if i != n-1 {
			stringified += ","
		}
	}
	stringified = "[" + stringified + "]"

	return stringified
}
