package logger

import (
	"strings"
)

type Options struct {
	rules   string
	url     string
	enabled interface{}
	queue   []string
}

const httpLoggerAgent string = "HttpLogger.go"

//base HttpLogger definition
type HttpLogger struct {
	*BaseLogger
	rules *HttpRules
}

// initialize HttpLogger
func NewHttpLogger(options Options) (*HttpLogger, error) {
	baseLogger := NewBaseLogger(httpLoggerAgent, options.url, options.enabled, options.queue)

	loggerRules, err := newHttpRules(options.rules)
	if err != nil {
		return nil, err
	}

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

	return logger, nil
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

	logger.Submit(msgStringify(details))
}
