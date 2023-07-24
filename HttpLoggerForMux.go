// © 2016-2023 Graylog, Inc.

package logger

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

var LIMIT = 1024 * 1024

type (
	// HttpLoggerForMux defines a struct used to log specifically gorilla/mux apps
	HttpLoggerForMux struct {
		HttpLogger *HttpLogger
		startTime  time.Time
		interval   time.Duration
		response   []byte
	}

	loggingResponseWriter struct { //custom response writer to wrap original writer in
		http.ResponseWriter
		loggingResp *http.Response
	}
)

// NewHttpLoggerForMux returns a pointer to an instance of an HttpLoggerForMux struct with the default options applied and an error.
// If there is no error, the error value returned will be nil.
func NewHttpLoggerForMux() (*HttpLoggerForMux, error) {

	HttpLogger, err := NewHttpLogger(Options{})

	if err != nil {
		return nil, err
	}

	httpLoggerForMux := HttpLoggerForMux{
		HttpLogger: HttpLogger,
		startTime:  time.Time{},
		interval:   0,
		response:   make([]byte, 0),
	}

	return &httpLoggerForMux, nil
}

// NewHttpLoggerForMuxOptions(o Options) returns a pointer to a HttpLoggerForMux struct with the given options o applied and an error.
// If there is no error, the error value returned will be nil.
func NewHttpLoggerForMuxOptions(options Options) (*HttpLoggerForMux, error) {

	HttpLogger, err := NewHttpLogger(options)

	if err != nil {
		return nil, err
	}

	httpLoggerForMux := HttpLoggerForMux{
		HttpLogger: HttpLogger,
		startTime:  time.Time{},
		interval:   0,
		response:   make([]byte, 0),
	}

	return &httpLoggerForMux, nil
}

// Write(b []byte) uses original response writer to write the body b to the client and then logs the response body.
// This is only used internally by response writer.
func (w *loggingResponseWriter) Write(body []byte) (int, error) { // uses original response writer to write and then logs the size

	size, err := w.ResponseWriter.Write(body)

	if err == nil {
		defer func() {
			w.loggingResp.Header = w.ResponseWriter.Header()
			if size > 0 {
				w.loggingResp.Header.Set("Content-Length", fmt.Sprint(size))
			}
		}()

		var loggedBodyBytes []byte
		if size < LIMIT {
			loggedBodyBytes = body
		} else {
			loggedBodyBytes = []byte(fmt.Sprintf("{ overflowed: %d }", size))
		}

		w.loggingResp = &http.Response{
			Header:     w.loggingResp.Header,
			Body:       io.NopCloser(bytes.NewBuffer(loggedBodyBytes)),
			StatusCode: w.loggingResp.StatusCode, // Status Code 200 will only be overridden if writeHeader is called
		}
	}

	return size, err
}

// WriteHeader(s int) uses original response writer to write the header with code s and then logs the response status code.
// This is only used internally by response writer.
func (w *loggingResponseWriter) WriteHeader(statusCode int) {
	w.loggingResp.StatusCode = statusCode

	w.ResponseWriter.WriteHeader(statusCode)
}

// LogData() takes 1 argument of type http.Handler and returns an object of the same type, http.Handler.
// This function is intended to be used in a Middleware function in a gorilla/mux server.
// For details on how to set up Middleware for a mux server see: https://github.com/resurfaceio/logger-go#logging_from_mux
func (muxLogger HttpLoggerForMux) LogData(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		loggingWriter := loggingResponseWriter{
			ResponseWriter: w,
			loggingResp: &http.Response{
				StatusCode: 200,
			},
		}

		buf, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		r.Body.Close()

		r.Body = io.NopCloser(bytes.NewBuffer(buf))

		loggingReq := &http.Request{
			Method:        r.Method,
			URL:           r.URL,
			Proto:         r.Proto,
			ProtoMajor:    r.ProtoMajor,
			ProtoMinor:    r.ProtoMinor,
			Header:        r.Header,
			ContentLength: r.ContentLength,
			Close:         r.Close,
			Host:          r.Host,
			Form:          r.Form,
			Trailer:       r.Trailer,
			RemoteAddr:    r.RemoteAddr,
			RequestURI:    r.RequestURI,
			TLS:           r.TLS,
			MultipartForm: r.MultipartForm,
			Response:      r.Response,
			Body:          io.NopCloser(bytes.NewBuffer(buf)),
		}

		now := time.Now()

		next.ServeHTTP(&loggingWriter, r)

		interval := time.Since(now).Milliseconds()

		SendHttpMessage(muxLogger.HttpLogger, loggingWriter.loggingResp, loggingReq, now.UnixNano()/int64(time.Millisecond), interval, nil)
	})
}
