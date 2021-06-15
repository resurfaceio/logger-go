package logger

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type HttpLoggerForMux struct {
	httpLogger HttpLogger
	router     mux.Router
}

func NewHttpLoggerForMux(r mux.Router) (*HttpLoggerForMux, error) {

	options := Options{}
	httpLogger, err := NewHttpLogger(options)

	if err != nil {
		return nil, err
	}

	httpLoggerForMux := HttpLoggerForMux{
		httpLogger: *httpLogger,
		router:     r,
	}

	return &httpLoggerForMux, nil
}

func NewHttpLoggerForMuxOptions(options Options, r mux.Router) (*HttpLoggerForMux, error) {

	httpLogger, err := NewHttpLogger(options)

	if err != nil {
		return nil, err
	}

	httpLoggerForMux := HttpLoggerForMux{
		httpLogger: *httpLogger,
		router:     r,
	}

	return &httpLoggerForMux, nil
}

func Log(next http.Handler) http.Handler { //WIP this is just to test middleware functionality
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Whale hello there!")

		next.ServeHTTP(w, r)
	})
}
