package logger

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type loggerOptions struct {
	rules string
	agent string
	url string
	enabled bool
	queue []string
}

const loggerAgent string = "HttpLogger.go"

//base HttpLogger definition
type HttpLogger struct{
	BaseLogger
	enabled bool
	queue []string
	skip_compression bool
	skip_submission bool
	rules string
	url string
}

// initialize HttpLogger 
func NewHttpLogger(options loggerOptions) *HttpLogger { 
	baseLogger := NewBaseLogger(options.agent, options.url, options.enabled, options.queue)

	logger := &HttpLogger{
		baseLogger,
	}

	logger.rules, err = NewHttpRules(_rules)
	
	logger.skipCompression = httpRules.skipCompression
	logger.skipSubmission = httpRules.skipSubmission

	if((logger.url != nil) && (strings.HasPrefix(logger.url, "http:") && !logger.rules.allowHttpUrl)) {
		logger.enableable = false;
		logger.enabled = false;
	}

	return logger
}

func (logger *HttpLgger) submitIfPassing(var details []string) {
	details = rules.apply(details) 

	if details == nil {
		return
	}

	details = append(["host", host]) // where does host come from?

	logger.submit(json.Unmarshal([]byte(details), &value)) 
}