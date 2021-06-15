package logger

import (
	"log"

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

func NewHttpLoggerForMuxOptions(options Options) (*HttpLoggerForMux, error) {

	httpLogger, err := NewHttpLogger(options)

	if err != nil {
		return nil, err
	}

	httpLoggerForMux := HttpLoggerForMux{
		httpLogger: *httpLogger,
		// router: r,
	}

	return &httpLoggerForMux, nil
}

func (loggerMux *HttpLoggerForMux) TestPrint() {
	log.Println("Whale hello there!")
}
