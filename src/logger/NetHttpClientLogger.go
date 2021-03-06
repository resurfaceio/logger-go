package logger

import (
	"io"
	"net/http"
	"net/url"
	"time"
)

type NetHttpClientLogger struct {
	http.Client
	httpLogger *HttpLogger
}

// construct new logger with given options struct{rules string, schema string}
func NewNetHttpClientLoggerOptions(options loggerOptions) *NetHttpClientLogger {
	return &NetHttpClientLogger{
		httpLogger: NewHttpLogger(options),
	}
}

// construct new logger without options
func NewNetHttpClientLogger() *NetHttpClientLogger {
	return &NetHttpClientLogger{
		httpLogger: NewHttpLogger(),
	}
}

func (logger *NetHttpClientLogger) Logger() *HttpLogger {
	return logger.httpLogger
}

// net.http.Client.CloseIdleConnections() wrapper
func (clientLogger *NetHttpClientLogger) CloseIdleConnections() {
	clientLogger.Client.CloseIdleConnections()
}

// net.http.Client.Do wrapper with logging
func (clientLogger *NetHttpClientLogger) Do(req *http.Request) (resp *http.Response, err error) {
	// start time for logging interval
	start := time.Now().UnixNano() / int64(time.Millisecond)

	logger := clientLogger.httpLogger

	// capture the response or error
	resp, err = clientLogger.Client.Do(req)

	if err != nil {
		return resp, err
	}

	// send logging message
	sendNetHttpClientMessage(logger, resp, start)

	return resp, err
}

// net.http.Client.Get wrapper with logging
func (clientLogger *NetHttpClientLogger) Get(url string) (resp *http.Response, err error) {
	// start time for logging interval
	start := time.Now().UnixNano() / int64(time.Millisecond)

	logger := clientLogger.httpLogger

	// capture the response or error
	resp, err = clientLogger.Client.Get(url)

	if err != nil {
		return resp, err
	}

	// send logging message
	sendNetHttpClientMessage(logger, resp, start)

	return resp, err
}

// net.http.Client.Head wrapper with logging
func (clientLogger *NetHttpClientLogger) Head(url string) (resp *http.Response, err error) {
	// start time for logging interval
	start := time.Now().UnixNano() / int64(time.Millisecond)

	logger := clientLogger.httpLogger

	// capture the response or error
	resp, err = clientLogger.Client.Head(url)

	if err != nil {
		return resp, err
	}

	// send logging message
	sendNetHttpClientMessage(logger, resp, start)

	return resp, err
}

// net.http.Client.Post wrapper with logging
func (clientLogger *NetHttpClientLogger) Post(url string, contentType string, body io.Reader) (resp *http.Response, err error) {
	// start time for logging interval
	start := time.Now().UnixNano() / int64(time.Millisecond)

	logger := clientLogger.httpLogger

	// capture the response or error
	resp, err = clientLogger.Client.Post(url, contentType, body)

	if err != nil {
		return resp, err
	}

	// send logging message
	sendNetHttpClientMessage(logger, resp, start)

	return resp, err
}

// net.http.Client.PostForm wrapper with logging
func (clientLogger *NetHttpClientLogger) PostForm(url string, data url.Values) (resp *http.Response, err error) {
	// start time for logging interval
	start := time.Now().UnixNano() / int64(time.Millisecond)

	logger := clientLogger.httpLogger

	// capture the response or error
	resp, err = clientLogger.Client.PostForm(url, data)

	if err != nil {
		return resp, err
	}

	// send logging message
	sendNetHttpClientMessage(logger, resp, start)

	return resp, err
}
