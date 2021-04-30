package logger

import (
	"encoding/json"
	"strings"
)

type Options struct {
	Rules   string
	Url     string
	Enabled interface{}
	Queue   []string
}

const httpLoggerAgent string = "HttpLogger.go"

//base HttpLogger definition
type HttpLogger struct {
	*BaseLogger
	rules *HttpRules
}

// initialize HttpLogger
func NewHttpLogger(options Options) (*HttpLogger, error) {
	baseLogger := NewBaseLogger(httpLoggerAgent, options.Url, options.Enabled, options.Queue)

	loggerRules, err := newHttpRules(options.Rules)
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

	byteStr, _ := json.Marshal(details)

	detailsString := string(byteStr)
	detailsString = strings.Replace(detailsString, "\\u003c", "<", -1)
	detailsString = strings.Replace(detailsString, "\\u003e", ">", -1)
	logger.Submit(detailsString)
}
