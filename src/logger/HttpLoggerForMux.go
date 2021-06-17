package logger

import (
	"log"
	"net/http"
	"time"
)

type HttpLoggerForMux struct {
	logger    HttpLogger
	startTime time.Time
	interval  time.Duration
	response  []byte
}

func NewHttpLoggerForMux() (*HttpLoggerForMux, error) {

	options := Options{}
	httpLogger, err := NewHttpLogger(options)

	if err != nil {
		return nil, err
	}

	httpLoggerForMux := HttpLoggerForMux{
		logger:    *httpLogger,
		startTime: time.Time{},
		interval:  0,
		response:  make([]byte, 0),
	}

	return &httpLoggerForMux, nil
}

func NewHttpLoggerForMuxOptions(options Options) (*HttpLoggerForMux, error) {

	httpLogger, err := NewHttpLogger(options)

	if err != nil {
		return nil, err
	}

	httpLoggerForMux := HttpLoggerForMux{
		logger:    *httpLogger,
		startTime: time.Time{},
		interval:  0,
		response:  make([]byte, 0),
	}

	return &httpLoggerForMux, nil
}

func (logger HttpLoggerForMux) StartResponse(next http.Handler) http.Handler { //WIP this is just to test middleware functionality
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// log.Println("Whale hello there!")
		logger.startTime = time.Now()

		next.ServeHTTP(w, r)

		logger.interval = time.Since(logger.startTime)
		log.Println("Start Time: ", logger.startTime, "Interval: ", logger.interval, "Method: ", r.Method, "Request Body: ", r.Body)
	})
}
