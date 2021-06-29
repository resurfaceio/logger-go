package logger

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"time"
)

type (
	HttpLoggerForMux struct {
		httpLogger *HttpLogger
		startTime  time.Time
		interval   time.Duration
		response   []byte
	}

	loggingResponseWriter struct { //custom response writer to wrap original writer in
		http.ResponseWriter
		loggingResp *http.Response
	}
)

func NewHttpLoggerForMux() (*HttpLoggerForMux, error) {

	options := Options{}
	httpLogger, err := NewHttpLogger(options)

	if err != nil {
		return nil, err
	}

	httpLoggerForMux := HttpLoggerForMux{
		httpLogger: httpLogger,
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
		httpLogger: httpLogger,
		startTime:  time.Time{},
		interval:   0,
		response:   make([]byte, 0),
	}

	return &httpLoggerForMux, nil
}

func (w *loggingResponseWriter) Write(body []byte) (int, error) { // uses original response writer to write and then logs the size

	w.loggingResp = &http.Response{
		Header:     w.loggingResp.Header,
		Body:       ioutil.NopCloser(bytes.NewBuffer(body)),
		StatusCode: w.loggingResp.StatusCode, // Status Code 200 will only be overriden if writeHeader is called
	}

	size, err := w.ResponseWriter.Write(body)

	// Copy header to log response
	w.loggingResp.Header = w.ResponseWriter.Header().Clone()

	// w.response.ContentLength += int64(size)
	return size, err
}

// uses original response writer to write the header and then logs the status code
func (w *loggingResponseWriter) WriteHeader(statusCode int) {
	w.loggingResp.StatusCode = statusCode

	w.ResponseWriter.WriteHeader(statusCode)
}

func (muxLogger HttpLoggerForMux) StartResponse(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// log.Println("Whale hello there!")

		customWriter := loggingResponseWriter{
			ResponseWriter: w,
			loggingResp: &http.Response{
				StatusCode: 200,
			},
		}

		next.ServeHTTP(&customWriter, r)

		muxLogger.startTime = time.Now()

		// log.Println(customWriter.loggingResp.Header)

		sendHttpMessage(muxLogger.httpLogger, customWriter.loggingResp, r, muxLogger.startTime)
	})
}
