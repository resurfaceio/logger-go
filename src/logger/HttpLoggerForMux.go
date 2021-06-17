package logger

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"time"
)

type (
	HttpLoggerForMux struct {
		httpLogger HttpLogger
		startTime  int64
		interval   time.Duration
		response   []byte
	}

	loggingResponseWriter struct { //custom response writer to wrap original writer in
		http.ResponseWriter
		response http.Response
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
		httpLogger: *httpLogger,
		startTime:  0,
		interval:   0,
		response:   make([]byte, 0),
	}

	return &httpLoggerForMux, nil
}

/*
body := "Hello world"
t := &http.Response{
  Status:        "200 OK",
  StatusCode:    200,
  Proto:         "HTTP/1.1",
  ProtoMajor:    1,
  ProtoMinor:    1,
  Body:          ioutil.NopCloser(bytes.NewBufferString(body)),
  ContentLength: int64(len(body)),
  Request:       req,
  Header:        make(http.Header, 0),
}
*/

func (w *loggingResponseWriter) Write(body []byte) (int, error) { // uses original response writer to write and then logs the size
	// w.response.Body = ioutil.NopCloser(bytes.NewBuffer(body)) // write body to response duplicate
	size, err := w.ResponseWriter.Write(body)
	// w.response.ContentLength += int64(size)
	return size, err
}

func (w *loggingResponseWriter) WriteHeader(code int) { // uses original response writer to write the header and then logs the status code
	// w.response.StatusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func (muxLogger HttpLoggerForMux) StartResponse(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// log.Println("Whale hello there!")

		body := "Hello world"
		resp := &http.Response{
			Status:        "200 OK",
			StatusCode:    200,
			Proto:         "HTTP/1.1",
			ProtoMajor:    1,
			ProtoMinor:    1,
			Body:          ioutil.NopCloser(bytes.NewBufferString(body)),
			ContentLength: int64(len(body)),
			Request:       nil,
			Header:        make(http.Header, 0),
		}

		customWriter := loggingResponseWriter{
			ResponseWriter: w,
			response:       *resp,
		}

		muxLogger.startTime = time.Now().UnixNano() / int64(time.Millisecond)

		sendHttpMessage(&muxLogger.httpLogger, resp, r, muxLogger.startTime)

		next.ServeHTTP(&customWriter, r) // replace standard response writer with custom one from above

	})
}
