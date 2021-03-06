package logger

import (
	"strings"
)

type loggerOptions struct {
	rules   string
	agent   string
	url     string
	enabled bool
	queue   []string
}

const loggerAgent string = "HttpLogger.go"

//base HttpLogger definition
type HttpLogger struct {
	*BaseLogger
	rules string
}

// initialize HttpLogger
func NewHttpLogger(options loggerOptions) *HttpLogger {
	baseLogger := NewBaseLogger(options.agent, options.url, options.enabled, options.queue)

	logger := &HttpLogger{
		baseLogger,
		"",
	}

	logger.rules, err = NewHttpRules(_rules)

	logger.skipCompression = httpRules.skipCompression
	logger.skipSubmission = httpRules.skipSubmission

	if (logger.url != "") && (strings.HasPrefix(logger.url, "http:") && !logger.rules.allowHttpUrl) {
		logger.enableable = false
		logger.enabled = false
	}

	return logger
}

func (logger *HttpLogger) submitIfPassing(details [][]string) {
	details = logger.rules.apply(details)

	if details == nil {
		return
	}

	details = append(details, []string{"host", logger.host})

	logger.submit(msgStringify(details))
}

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
