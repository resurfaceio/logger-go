// © 2016-2022 Resurface Labs Inc.

package logger

import (
	"io"
	"net/http"
	"net/url"
	"time"
)

//NetHttpClientLogger defines a struct used to log specifically from the client side of API interactions using the net/http package.
type NetHttpClientLogger struct {
	http.Client
	HttpLogger *HttpLogger
}

//NewNetHttpClientLoggerOptions() takes 1 argument of type logger.Options and returns 2 objects; a pointer to an instance of an NetHttpClientLogger struct and an error.
//The NetHttpClientLogger returned by this function has the given options applied.
//If there is no error, the error value returned will be nil.
func NewNetHttpClientLoggerOptions(options Options) (*NetHttpClientLogger, error) {
	HttpLogger, err := NewHttpLogger(options)
	if err != nil {
		return nil, err
	}
	return &NetHttpClientLogger{
		HttpLogger: HttpLogger,
	}, nil
}

//NewNetHttpClientLogger() takes no arguments and returns 2 objects; a pointer to an instance of an NetHttpClientLogger struct and an error.
//The NetHttpClientLogger returned by this function has the default options applied.
//If there is no error, the error value returned will be nil.
func NewNetHttpClientLogger() (*NetHttpClientLogger, error) {
	options := Options{}
	HttpLogger, err := NewHttpLogger(options)
	if err != nil {
		return nil, err
	}
	return &NetHttpClientLogger{
		HttpLogger: HttpLogger,
	}, nil
}

func (logger *NetHttpClientLogger) Logger() *HttpLogger {
	return logger.HttpLogger
}

//net.http.Client.CloseIdleConnections() wrapper
func (clientLogger *NetHttpClientLogger) CloseIdleConnections() {
	clientLogger.Client.CloseIdleConnections()
}

//net.http.Client.Do wrapper with logging
func (clientLogger *NetHttpClientLogger) Do(req *http.Request) (resp *http.Response, err error) {
	// start time for logging interval
	logger := clientLogger.HttpLogger

	now := time.Now()

	// capture the response or error
	resp, err = clientLogger.Client.Do(req)

	interval := time.Since(now).Milliseconds()

	if err != nil {
		return resp, err
	}

	// send logging message
	SendHttpMessage(logger, resp, resp.Request, (now.Unix() * int64(time.Millisecond)), interval, nil)

	return resp, err
}

// net.http.Client.Get wrapper with logging
func (clientLogger *NetHttpClientLogger) Get(url string) (resp *http.Response, err error) {
	// start time for logging interval
	logger := clientLogger.HttpLogger

	now := time.Now()

	// capture the response or error
	// Devin 03/31/2021
	// Something happens here with the .Get where err does return an error = "unsupported protocol scheme"
	resp, err = clientLogger.Client.Get(url)

	interval := time.Since(now).Milliseconds()

	if err != nil {
		return resp, err
	}

	// send logging message
	SendHttpMessage(logger, resp, resp.Request, (now.Unix() * int64(time.Millisecond)), interval, nil)

	return resp, err
}

// net.http.Client.Head wrapper with logging
func (clientLogger *NetHttpClientLogger) Head(url string) (resp *http.Response, err error) {
	// start time for logging interval

	logger := clientLogger.HttpLogger

	now := time.Now()

	// capture the response or error
	resp, err = clientLogger.Client.Head(url)

	interval := time.Since(now).Milliseconds()

	if err != nil {
		return resp, err
	}

	// send logging message
	SendHttpMessage(logger, resp, resp.Request, (now.Unix() * int64(time.Millisecond)), interval, nil)

	return resp, err
}

// net.http.Client.Post wrapper with logging
func (clientLogger *NetHttpClientLogger) Post(url string, contentType string, body io.Reader) (resp *http.Response, err error) {
	// start time for logging interval

	logger := clientLogger.HttpLogger

	now := time.Now()

	// capture the response or error
	resp, err = clientLogger.Client.Post(url, contentType, body)

	interval := time.Since(now).Milliseconds()

	if err != nil {
		return resp, err
	}

	// send logging message
	SendHttpMessage(logger, resp, resp.Request, (now.Unix() * int64(time.Millisecond)), interval, nil)

	return resp, err
}

// net.http.Client.PostForm wrapper with logging
func (clientLogger *NetHttpClientLogger) PostForm(url string, data url.Values) (resp *http.Response, err error) {
	// start time for logging interval

	logger := clientLogger.HttpLogger

	now := time.Now()

	// capture the response or error
	resp, err = clientLogger.Client.PostForm(url, data)

	interval := time.Since(now).Milliseconds()

	if err != nil {
		return resp, err
	}

	// send logging message
	SendHttpMessage(logger, resp, resp.Request, (now.Unix() * int64(time.Millisecond)), interval, nil)

	return resp, err
}
