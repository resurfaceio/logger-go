// © 2016-2021 Resurface Labs Inc.

//Resurface Go Logger provides tools to log API requests and responses from different Golang web frameworks to a complete API system of record. (https://resurface.io)
package logger

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
)

//Options struct is passed to a "NewLogger" function to specifiy the desired configuration of the logger to be created.
type Options struct {
	//Rules defines the rules that will be applied to the logger.
	Rules string

	//Url defines the Url the logger will send the logs to.
	Url string

	//Enabled defines the state of the logger; enabled or disabled.
	Enabled interface{}

	//Queue is a slice of strings used to store logs; exclusively for testing purposes.
	//Queue must be nil for the logger to properly function.
	Queue []string
}

const httpLoggerAgent string = "HttpLogger.go"

//base httpLogger definition
type httpLogger struct {
	*baseLogger
	rules *HttpRules
}

// NewHttpLogger returns a pointer to a new httpLogger object, with the given options applied, and an error
func NewHttpLogger(options Options) (*httpLogger, error) {
	baseLogger := newBaseLogger(httpLoggerAgent, options.Url, options.Enabled, options.Queue)

	loggerRules, err := newHttpRules(options.Rules)
	if err != nil {
		return nil, err
	}

	logger := &httpLogger{
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

func (logger *httpLogger) submitIfPassing(msg [][]string) {
	msg = logger.rules.apply(msg)

	if msg == nil {
		return
	}

	msg = append(msg, []string{"host", logger.host})

	byteStr, _ := json.Marshal(msg)

	msgString := string(byteStr)
	msgString = strings.Replace(msgString, "\\u003c", "<", -1)
	msgString = strings.Replace(msgString, "\\u003e", ">", -1)
	logger.submit(msgString)
}

// global client to avoid opening a new connection for every request
var httpLoggerClient *http.Client

func init() {
	tr := &http.Transport{
		MaxIdleConnsPerHost: 10000,
		TLSHandshakeTimeout: 0 * time.Second,
	}
	httpLoggerClient = &http.Client{Transport: tr}

	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
