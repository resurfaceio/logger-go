package logger

import (
	"bytes"
	"io/ioutil"
	"log"
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
		httpResp *http.Response
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
	w.httpResp = &http.Response{
		Body: ioutil.NopCloser(bytes.NewBuffer(body)), // write body to response duplicate,
	}

	size, err := w.ResponseWriter.Write(body)
	// w.response.ContentLength += int64(size)
	return size, err
}

func (w *loggingResponseWriter) WriteHeader(i int) { // uses original response writer to write the header and then logs the status code
	// w.httpResp.Header = w.Header()
	log.Println("Status Code: ", i)
	// w.httpResp.StatusCode = i
	w.ResponseWriter.WriteHeader(i)
}

func (muxLogger HttpLoggerForMux) StartResponse(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// log.Println("Whale hello there!")

		//test response for logger
		// body := "Hello world"
		// resp := &http.Response{
		// 	Status:        "200 OK",
		// 	StatusCode:    200,
		// 	Proto:         "HTTP/1.1",
		// 	ProtoMajor:    1,
		// 	ProtoMinor:    1,
		// 	Body:          ioutil.NopCloser(bytes.NewBufferString(body)),
		// 	ContentLength: int64(len(body)),
		// 	Request:       nil,
		// 	Header:        make(http.Header),
		// }

		customWriter := loggingResponseWriter{
			ResponseWriter: w,
			httpResp:       &http.Response{},
		}

		// muxLogger.startTime = time.Now().UnixNano() / int64(time.Millisecond)

		// sendHttpMessage(muxLogger.httpLogger, resp, r, muxLogger.startTime)

		next.ServeHTTP(&customWriter, r) // replace standard response writer with custom one from above

		log.Println("Header - User-Agent: ", customWriter.httpResp.Header.Get("User-Agent"))

	})
}
