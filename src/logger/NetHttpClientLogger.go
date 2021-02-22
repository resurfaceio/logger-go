package logger

import (
	"net/http"
	"time"
)

type NetHttpClientLogger struct {
	http.Client
	httpLogger *HttpLogger
}

func NewNetHttpClientLogger(options interface{}) *NetHttpClientLogger {
	return &NetHttpClientLogger{
		httpLogger: NewHttpLogger(options),
	}
}

func (logger *NetHttpClientLogger) Logger() *HttpLogger {
	return logger.httpLogger
}

func (clientLogger *NetHttpClientLogger) Get(url string) (resp *http.Response, err error) {
	// start time for logging interval
	start := time.Now().UnixNano() / int64(time.Millisecond)

	logger := clientLogger.httpLogger

	// capture the response or error
	resp, err = clientLogger.Client.Get(url)

	if err != nil {
		return resp, err
	}

	// before sending should we first check "if (status < 300 || status === 302)"?

	// now = time.Now().UnixNano() / int64(time.Millisecond)
	sendNetHttpClientMessage(logger, resp, start)

	return resp, err
}
