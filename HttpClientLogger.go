// © 2016-2021 Resurface Labs Inc.

package logger

import (
	"io"
	"net/http"
	"net/url"
	"time"
)

type NetHttpClientLogger struct {
	http.Client
	httpLogger *httpLogger
}

//NewNetHttpClientLoggerOptions() takes 1 argument of type logger.Options and returns 2 objects; a pointer to an instance of an NetHttpClientLogger struct and an error.
//The NetHttpClientLogger returned by this function has the given options applied.
//If there is no error, the error value returned will be nil.
func NewNetHttpClientLoggerOptions(options Options) (*NetHttpClientLogger, error) {
	httpLogger, err := newHttpLogger(options)
	if err != nil {
		return nil, err
	}
	return &NetHttpClientLogger{
		httpLogger: httpLogger,
	}, nil
}

//NewNetHttpClientLogger() takes no arguments and returns 2 objects; a pointer to an instance of an NetHttpClientLogger struct and an error.
//The NetHttpClientLogger returned by this function has the default options applied.
//If there is no error, the error value returned will be nil.
func NewNetHttpClientLogger() (*NetHttpClientLogger, error) {
	options := Options{}
	httpLogger, err := newHttpLogger(options)
	if err != nil {
		return nil, err
	}
	return &NetHttpClientLogger{
		httpLogger: httpLogger,
	}, nil
}

func (logger *NetHttpClientLogger) Logger() *httpLogger {
	return logger.httpLogger
}

// net.http.Client.CloseIdleConnections() wrapper
func (clientLogger *NetHttpClientLogger) CloseIdleConnections() {
	clientLogger.Client.CloseIdleConnections()
}

// net.http.Client.Do wrapper with logging
func (clientLogger *NetHttpClientLogger) Do(req *http.Request) (resp *http.Response, err error) {
	// start time for logging interval
	start := time.Now()

	logger := clientLogger.httpLogger

	// capture the response or error
	resp, err = clientLogger.Client.Do(req)

	if err != nil {
		return resp, err
	}

	// send logging message
	sendHttpMessage(logger, resp, resp.Request, start)

	return resp, err
}

// net.http.Client.Get wrapper with logging
func (clientLogger *NetHttpClientLogger) Get(url string) (resp *http.Response, err error) {
	// start time for logging interval
	start := time.Now()

	logger := clientLogger.httpLogger

	// capture the response or error
	// Devin 03/31/2021
	// Something happens here with the .Get where err does return an error = "unsupported protocol scheme"
	resp, err = clientLogger.Client.Get(url)

	if err != nil {
		return resp, err
	}

	// send logging message
	sendHttpMessage(logger, resp, resp.Request, start)

	return resp, err
}

// net.http.Client.Head wrapper with logging
func (clientLogger *NetHttpClientLogger) Head(url string) (resp *http.Response, err error) {
	// start time for logging interval
	start := time.Now()

	logger := clientLogger.httpLogger

	// capture the response or error
	resp, err = clientLogger.Client.Head(url)

	if err != nil {
		return resp, err
	}

	// send logging message
	sendHttpMessage(logger, resp, resp.Request, start)

	return resp, err
}

// net.http.Client.Post wrapper with logging
func (clientLogger *NetHttpClientLogger) Post(url string, contentType string, body io.Reader) (resp *http.Response, err error) {
	// start time for logging interval
	start := time.Now()

	logger := clientLogger.httpLogger

	// capture the response or error
	resp, err = clientLogger.Client.Post(url, contentType, body)

	if err != nil {
		return resp, err
	}

	// send logging message
	sendHttpMessage(logger, resp, resp.Request, start)

	return resp, err
}

// net.http.Client.PostForm wrapper with logging
func (clientLogger *NetHttpClientLogger) PostForm(url string, data url.Values) (resp *http.Response, err error) {
	// start time for logging interval
	start := time.Now()

	logger := clientLogger.httpLogger

	// capture the response or error
	resp, err = clientLogger.Client.PostForm(url, data)

	if err != nil {
		return resp, err
	}

	// send logging message
	sendHttpMessage(logger, resp, resp.Request, start)

	return resp, err
}
