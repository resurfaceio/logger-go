package logger

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
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

	defer func() {
		w.loggingResp.Header = w.ResponseWriter.Header()
		if len(body) > 0 {
			w.loggingResp.Header.Set("Content-Length", fmt.Sprint(len(body)))
		}
	}()

	return size, err
}

// uses original response writer to write the header and then logs the status code
func (w *loggingResponseWriter) WriteHeader(statusCode int) {
	w.loggingResp.StatusCode = statusCode

	w.ResponseWriter.WriteHeader(statusCode)
}

func (muxLogger HttpLoggerForMux) StartResponse(next http.Handler) http.Handler {
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

		r.Body = ioutil.NopCloser(bytes.NewBuffer(buf))

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
			Body:          ioutil.NopCloser(bytes.NewBuffer(buf)),
		}

		next.ServeHTTP(&loggingWriter, r)
		// loggingWriter.loggingResp.Header = loggingWriter.ResponseWriter.Header()
		log.Println(loggingWriter.loggingResp.Header)

		muxLogger.startTime = time.Now()

		sendHttpMessage(muxLogger.httpLogger, loggingWriter.loggingResp, loggingReq, muxLogger.startTime)
	})
}
