package logger

import (
	"log"
	"net/http"
	"time"
)

type (
	HttpLoggerForMux struct {
		httpLogger HttpLogger
		startTime  time.Time
		interval   time.Duration
		response   []byte
	}

	responseData struct { //to hold response data
		status int
		body   string
		size   int
	}

	loggingResponseWriter struct { //custom response writer to wrap original writer in
		http.ResponseWriter
		responseData *responseData
	}
)

func NewHttpLoggerForMux() (*HttpLoggerForMux, error) {

	options := Options{}
	httpLogger, err := NewHttpLogger(options)

	if err != nil {
		return nil, err
	}

	httpLoggerForMux := HttpLoggerForMux{
		httpLogger: *httpLogger,
		startTime:  time.Time{},
		interval:   0,
		response:   make([]byte, 0),
	}

	return &httpLoggerForMux, nil
}

func NewHttpLoggerForMuxOptions(options Options) (*HttpLoggerForMux, error) {

	httpLogger, err := NewHttpLogger(options)

	if err != nil {
		return nil, err
	}

	httpLoggerForMux := HttpLoggerForMux{
		httpLogger: *httpLogger,
		startTime:  time.Time{},
		interval:   0,
		response:   make([]byte, 0),
	}

	return &httpLoggerForMux, nil
}

func (w *loggingResponseWriter) Write(body []byte) (int, error) { // uses original response writer to write and then logs the size
	w.responseData.body = string(body)
	size, err := w.ResponseWriter.Write(body)
	w.responseData.size += size
	return size, err
}

func (w *loggingResponseWriter) WriteHeader(code int) { // uses original response writer to write the header and then logs the status code
	w.responseData.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (muxLogger HttpLoggerForMux) StartResponse(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// log.Println("Whale hello there!")

		responseData := &responseData{
			status: 0,
			body:   "",
			size:   0,
		}

		customWriter := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}

		muxLogger.startTime = time.Now()

		next.ServeHTTP(&customWriter, r) // replace standard response writer with custom one from above

		muxLogger.interval = time.Since(muxLogger.startTime)
		log.Println("Response Status: ", responseData.status, " Response Body: ", responseData.body, " Interval: ", muxLogger.interval, " Method: ", r.Method, " Request Body: ", r.Body)
	})
}
