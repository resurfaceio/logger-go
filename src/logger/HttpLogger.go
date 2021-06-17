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

func (logger *HttpLogger) submitIfPassing(msg [][]string) {
	msg = logger.rules.apply(msg)

	if msg == nil {
		return
	}

	msg = append(msg, []string{"host", logger.host})

	byteStr, _ := json.Marshal(msg)

	msgString := string(byteStr)
	msgString = strings.Replace(msgString, "\\u003c", "<", -1)
	msgString = strings.Replace(msgString, "\\u003e", ">", -1)
	logger.Submit(msgString)
}
