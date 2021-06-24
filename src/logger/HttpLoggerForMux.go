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
		startTime  int64
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
		startTime:  0,
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
		startTime:  0,
		interval:   0,
		response:   make([]byte, 0),
	}

	return &httpLoggerForMux, nil
}

func (w *loggingResponseWriter) Write(body []byte) (int, error) { // uses original response writer to write and then logs the size
	w.loggingResp = &http.Response{
		Body: ioutil.NopCloser(bytes.NewBuffer(body)), // write body to response duplicate,
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
				// initialize status code to 200 incase WriteHeader is not called
				StatusCode: 200,
			},
		}

		// replace standard response writer with custom one from above
		next.ServeHTTP(&customWriter, r)
		// log.Println(customWriter.httpResp.StatusCode)

		muxLogger.startTime = time.Now().UnixNano() / int64(time.Millisecond)

		sendHttpMessage(muxLogger.httpLogger, customWriter.loggingResp, r, muxLogger.startTime)
	})
}
